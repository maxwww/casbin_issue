package main

import (
	"log"

	"github.com/casbin/casbin/v2"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	a := fileadapter.NewAdapter("./policy.csv")

	e, err := casbin.NewSyncedEnforcer("./model.conf", a)
	if err != nil {
		panic(err)
	}

	e.EnableAutoSave(true)

	app.Get("/list", listPolicies(e))
	app.Get("/assign", assignPolicy(e))
	app.Get("/check", check(e))

	log.Println(app.Listen(":3000"))
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
