package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"goly/model"
	"goly/utils"
	"strconv"
)

func getAllRedirects(ctx *fiber.Ctx) error {
	golies, err := model.GetAllGolies()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error getting all goly links",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(golies)
}

func getGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could not parse id",
		})
	}
	goly, err := model.GetGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could not retreive goly from db",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func createGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("applications/json")

	var goly model.Goly

	err := ctx.BodyParser(&goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON",
		})
	}

	if goly.Random {
		goly.Goly, err = utils.NanoId(16)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "InternalServerError",
			})
		}
	}
	err = model.CreateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cloud not create goly in db",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(goly)
}

func updateGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("applications/json")
	var goly model.Goly
	err := ctx.BodyParser(&goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON",
		})
	}
	err = model.UpdateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cloud not update goly link in DB",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func deleteGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could not parse id",
		})
	}
	err = model.DeleteGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could delete from DB",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "goly deleted !",
	})
}

func redirect(ctx *fiber.Ctx) error {
	golyURL := ctx.Params("redirect")

	goly, err := model.FindByGolyUrl(golyURL)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "goly not found",
		})
	}

	goly.Clicked += 1

	err = model.UpdateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "cloud not update goly link in DB",
		})
	}

	return ctx.Redirect(goly.Redirect, fiber.StatusTemporaryRedirect)
}

func SetupAndListen() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	router.Get("/r/:redirect", redirect)

	router.Get("/goly/", getAllRedirects)
	router.Get("/goly/:id", getGoly)
	router.Post("/goly", createGoly)
	router.Patch("/goly", updateGoly)
	router.Delete("/goly/:id", deleteGoly)
	router.Listen(":3000")
}
