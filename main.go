package main

import (
	"database/sql"
	"log"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/test?parseTime=true&time_zone=UTC")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	casbinModel, err := model.NewModelFromString(CasbinModel)
	if err != nil {
		panic(err)
	}

	a, err := sqladapter.NewAdapter(db, "mysql", "casbin_rule")
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewSyncedEnforcer(casbinModel, a)
	if err != nil {
		panic(err)
	}

	e.EnableAutoSave(true)

	app.Get("/init", initPolicies(db, e))
	app.Get("/list", listPolicies(e))
	app.Get("/assign", assignPolicy(e))
	app.Get("/check", check(e))

	log.Println(app.Listen(":3000"))
}

func initPolicies(db *sql.DB, e *casbin.SyncedEnforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := db.Exec("delete from casbin_rule;")
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = e.LoadPolicy()

		policies := [][]any{
			{"root:root:root", "*", "*"},
			{"core:users:list", "/users", "GET"},
			{"core:users:create", "/users", "POST"},
		}

		for _, p := range policies {
			_, err := e.AddPolicy(p...)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		groups := [][]any{
			// Root role
			{"Root", "root:root:root"},
			// List Users Roles
			{"List Users Roles", "core:users:list"},
			// Admin Users Roles
			{"Admin Users Roles", "core:users:list"},
			{"Admin Users Roles", "core:users:create"},
			// Root user
			{"1", "Root"},
			// List Users Roles user
			{"2", "List Users Roles"},
			// Admin Users Roles users
			{"3", "Admin Users Roles"},
			{"4", "Admin Users Roles"},
		}

		for _, g := range groups {
			_, err := e.AddGroupingPolicy(g...)
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		err = e.LoadPolicy()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(map[string]string{
			"status": "ok",
		})
	}
}

func listPolicies(e *casbin.SyncedEnforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		groupingPolicy := e.GetGroupingPolicy()
		return c.JSON(groupingPolicy)
	}
}

func assignPolicy(e *casbin.SyncedEnforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type AssignRole struct {
			UserId string `query:"user_id"`
			Role   string `query:"role"`
		}
		var assignRole AssignRole
		if err := c.QueryParser(&assignRole); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		log.Println(assignRole)

		permissions, err := e.GetRolesForUser(assignRole.Role)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if len(permissions) == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, err = e.DeleteRolesForUser(assignRole.UserId)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		_, err = e.AddRoleForUser(assignRole.UserId, assignRole.Role)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func check(e *casbin.SyncedEnforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type User struct {
			UserId string `query:"user_id"`
		}
		var user User
		if err := c.QueryParser(&user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		log.Println(user)

		allowed, err := e.Enforce(user.UserId, "/users", "PUT")
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(map[string]any{
			"allowed": allowed,
		})
	}
}
