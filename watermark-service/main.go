package main

import (
        "bytes"
        "image"
        "image/draw"
        "image/jpeg"
        _ "image/png"
        "log"

        "github.com/gofiber/fiber/v2"
        "github.com/nfnt/resize"
)

func main() {
        app := fiber.New(fiber.Config{
                DisableStartupMessage: false,
        })

        app.Post("/sellar", func(c *fiber.Ctx) error {
                log.Println("--> Petición de sellado recibida desde n8n")

                // 1. Recibir la foto de n8n
                fotoHeader, errFoto := c.FormFile("foto")
                if errFoto != nil {
                        log.Println("ERROR: Falta el campo 'foto' en el form-data")
                        return c.Status(400).SendString("Falta la foto: " + errFoto.Error())
                }
                fotoFile, errFoto := fotoHeader.Open()
                if errFoto != nil {
                        return c.Status(500).SendString(errFoto.Error())
                }
                defer fotoFile.Close()

                // 2. Recibir el logo de n8n
                logoHeader, errLogo := c.FormFile("logo")
                if errLogo != nil {
                        log.Println("ERROR: Falta el campo 'logo' en el form-data")
                        return c.Status(400).SendString("Falta el logo: " + errLogo.Error())
                }
                logoFile, errLogo := logoHeader.Open()
                if errLogo != nil {
                        return c.Status(500).SendString(errLogo.Error())
                }
                defer logoFile.Close()

                // 3. Decodificar las imágenes
                imgFoto, _, errDec := image.Decode(fotoFile)
                if errDec != nil {
                        log.Println("ERROR: No se pudo decodificar la foto base")
                        return c.Status(500).SendString("Error al decodificar foto: " + errDec.Error())
                }
                imgLogo, _, errDec := image.Decode(logoFile)
                if errDec != nil {
                        log.Println("ERROR: No se pudo decodificar el logo")
                        return c.Status(500).SendString("Error al decodificar logo: " + errDec.Error())
                }

                // 4. Calcular redimensión del logo (25% del ancho de la foto)
                boundsFoto := imgFoto.Bounds()
                anchoFoto := boundsFoto.Dx()
                nuevoAnchoLogo := uint(anchoFoto) / 4

                logoRedimensionado := resize.Resize(nuevoAnchoLogo, 0, imgLogo, resize.Lanczos3)
                boundsLogo := logoRedimensionado.Bounds()

                // 5. Crear el lienzo RGBA
                rgba := image.NewRGBA(boundsFoto)
                draw.Draw(rgba, boundsFoto, imgFoto, image.Point{}, draw.Src)

                // Calcular coordenadas
                posX := boundsFoto.Max.X - boundsLogo.Max.X - 40
                posY := boundsFoto.Max.Y - boundsLogo.Max.Y - 40
                offset := image.Pt(posX, posY)

                // Superponer el logo
                draw.Draw(rgba, boundsLogo.Add(offset), logoRedimensionado, image.Point{}, draw.Over)

                // 6. Codificar el resultado final a JPEG
                var buf bytes.Buffer
                errEnc := jpeg.Encode(&buf, rgba, &jpeg.Options{Quality: 90})
                if errEnc != nil {
                        log.Println("ERROR: Fallo crítico al codificar la salida JPEG")
                        return c.Status(500).SendString("Error al codificar resultado: " + errEnc.Error())
                }

                log.Println("--> ¡Imagen procesada con éxito en Go!")
                c.Set("Content-Type", "image/jpeg")
                return c.Send(buf.Bytes())
        })

        log.Fatal(app.Listen(":5000"))
