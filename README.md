<h3 align="center">
  <img src="logo.png" width="500" alt="Logo"/><br/>
</h3>

GoMiniWin es un mini-conjunto de funciones para abrir una ventana, pintar en ella y detectar la presión de algunas teclas. Lo justo para poder implementar juegos sencillos sin necesidad de conocer el API de Windows o Linux.

## Inpiración

Esta librería es un pseudo-port de la librería [MiniWin](https://github.com/pauek/MiniWin) realizada por Pauek pero para Go (la original es para c++). Créditos a él y échenle un ojito a su [canal](https://www.youtube.com/@pauek), es un grande!

## Instalación

Asegúrate de estar en Windows o Linux.

```go
go get github.com/jibaru/gominiwin
```

## Plantilla básica

Para iniciar una nueva ventana, deberás importar el paquete necesario de acuerdo al sistema operativo que estés usando.

Aquí una plantilla para windows (para linux, deberás usar "github.com/jibaru/gominiwin/linux" en lugar de "github.com/jibaru/gominiwin/windows").

```go
package main

import (
	"time"

	"github.com/jibaru/gominiwin/colors"
	"github.com/jibaru/gominiwin/keys"
	"github.com/jibaru/gominiwin/windows"
)

func main() {
    // New incializa una nueva ventana, en este caso de 800x600
	w, err := windows.New("Example", 800, 600)
	if err != nil {
		panic(err)
	}

    // Podemos pintar en la ventana usando una gorutina y un bucle infinito
	go func() {
		for {
		    // Pinta "Hola GoMiniwin!" en la posición (100, 100)
		    w.SetText(100, 100, "Hola GoMiniwin!")
			// Muestra lo pintado
			w.Refresh()
			// Tiempo de espera para no repintar tan rápido
			time.Sleep(36 * time.Millisecond)
		}
	}()

    // Start inicializa la ventana, esta bloquea la gorutina principal
	w.Start()
}
```

## Funcionamiento de la ventana

Al ejecutar `New(...)`, se crea una nueva ventana. Actualmente GoMiniwin solo soporta el uso de una sola ventana, asi que evita crear múltiples ventanas. Esta ventana no se puede redimensionar con el ratón, pero puedes hacer con la función `Resize`.

Algo importante a tener en cuenta es que las coordenadas de la ventana son las coordenadas de una "matriz de datos", donde el X crece a medida que vas a la derecha, y la Y crece a medida que vas hacia abajo. El punto (0,0) está en la esquina superior izquierda.

## Funciones

### (windows/linux).New

Crea una nueva ventana.

| Argumentos | Descripción             |
| ---------- | ----------------------- |
| `title`    | El título de la ventana |
| `width`    | Ancho de la ventana     |
| `height`   | Alto de la ventana      |

| Valores de retorno | Descripción                                                                                                         |
| ------------------ | ------------------------------------------------------------------------------------------------------------------- |
| \*windows.window   | Puntero a la ventana. Esta contiene todas las operaciones para pintar o redimensionar. Implementa la interfaz `Win` |
| error              | En caso de error, mensaje informativo                                                                               |

## Documentación de la interfaz `Win`

#### `Point(x, y float32)`

Dibuja un punto en las coordenadas especificadas (usa el color actual).

| Argumentos | Descripción             |
| ---------- | ----------------------- |
| `x`        | Coordenada X del punto. |
| `y`        | Coordenada Y del punto. |

#### `Line(x1, y1, x2, y2 float32)`

Dibuja una línea desde el punto `(x1, y1)` hasta el punto `(x2, y2)` (usa el color actual).

| Argumentos | Descripción                     |
| ---------- | ------------------------------- |
| `x1`       | Coordenada X del punto inicial. |
| `y1`       | Coordenada Y del punto inicial. |
| `x2`       | Coordenada X del punto final.   |
| `y2`       | Coordenada Y del punto final.   |

#### `Rectangle(left, top, right, bottom float32)`

Dibuja un rectángulo definido por las esquinas superior izquierda (`left`, `top`) e inferior derecha (`right`, `bottom`). (usa el color actual para el borde).

| Argumentos | Descripción                      |
| ---------- | -------------------------------- |
| `left`     | Coordenada X del lado izquierdo. |
| `top`      | Coordenada Y del lado superior.  |
| `right`    | Coordenada X del lado derecho.   |
| `bottom`   | Coordenada Y del lado inferior.  |

#### `FilledRectangle(left, top, right, bottom float32)`

Dibuja un rectángulo relleno (usa el color actual).

| Argumentos | Descripción                      |
| ---------- | -------------------------------- |
| `left`     | Coordenada X del lado izquierdo. |
| `top`      | Coordenada Y del lado superior.  |
| `right`    | Coordenada X del lado derecho.   |
| `bottom`   | Coordenada Y del lado inferior.  |

#### `Circle(centerX, centerY, radius float32)`

Dibuja un círculo (usa el color actual para el borde).

| Argumentos | Descripción                          |
| ---------- | ------------------------------------ |
| `centerX`  | Coordenada X del centro del círculo. |
| `centerY`  | Coordenada Y del centro del círculo. |
| `radius`   | Radio del círculo.                   |

#### `FilledCircle(centerX, centerY, radius float32)`

Dibuja un círculo relleno (usa el color actual).

| Argumentos | Descripción                          |
| ---------- | ------------------------------------ |
| `centerX`  | Coordenada X del centro del círculo. |
| `centerY`  | Coordenada Y del centro del círculo. |
| `radius`   | Radio del círculo.                   |

#### `SetColor(c colors.Color)`

Establece el color actual para las operaciones de dibujo.

| Argumentos | Descripción                      |
| ---------- | -------------------------------- |
| `c`        | Color en formato `colors.Color`. |

#### `SetColorRGB(r, g, b int)`

Establece el color actual utilizando componentes RGB.

| Argumentos | Descripción               |
| ---------- | ------------------------- |
| `r`        | Componente rojo (0-255).  |
| `g`        | Componente verde (0-255). |
| `b`        | Componente azul (0-255).  |

#### `SetText(x, y float32, content string)`

Dibuja un texto en las coordenadas especificadas.

| Argumentos | Descripción             |
| ---------- | ----------------------- |
| `x`        | Coordenada X del texto. |
| `y`        | Coordenada Y del texto. |
| `content`  | Texto a dibujar.        |

#### `KeyPressed() int`

Devuelve el código de la última tecla presionada.

| Valores de retorno | Descripción                                |
| ------------------ | ------------------------------------------ |
| `int`              | Código de la tecla presionada (si existe). |

#### `MouseState() (bool, float32, float32)`

Obtiene el estado actual del mouse.

| Valores de retorno | Descripción                   |
| ------------------ | ----------------------------- |
| `bool`             | Indica si hay un clic activo. |
| `float32`          | Coordenada X del mouse.       |
| `float32`          | Coordenada Y del mouse.       |

#### `IsMouseInside() bool`

Devuelve si el mouse está dentro de la ventana.

| Valores de retorno | Descripción                                                |
| ------------------ | ---------------------------------------------------------- |
| `bool`             | `true` si el mouse está dentro, `false` en caso contrario. |

#### `MouseX() float32`

Devuelve la coordenada X actual del mouse.

| Valores de retorno | Descripción             |
| ------------------ | ----------------------- |
| `float32`          | Coordenada X del mouse. |

#### `MouseY() float32`

Devuelve la coordenada Y actual del mouse.

| Valores de retorno | Descripción             |
| ------------------ | ----------------------- |
| `float32`          | Coordenada Y del mouse. |

#### `MouseButtons() (bool, bool)`

Devuelve el estado de los botones del mouse.

| Valores de retorno | Descripción                                    |
| ------------------ | ---------------------------------------------- |
| `bool`             | Estado del botón izquierdo (`true` o `false`). |
| `bool`             | Estado del botón derecho (`true` o `false`).   |

#### `MouseLeftClicked() bool`

Devuelve si el botón izquierdo del mouse fue presionado.

| Valores de retorno | Descripción                        |
| ------------------ | ---------------------------------- |
| `bool`             | `true` si el botón fue presionado. |

#### `MouseRightClicked() bool`

Devuelve si el botón derecho del mouse fue presionado.

| Valores de retorno | Descripción                        |
| ------------------ | ---------------------------------- |
| `bool`             | `true` si el botón fue presionado. |

#### `Start()`

Inicia el bucle principal de la ventana. Este bloquea la gorutina principal.

#### `Clear()`

Limpia el contenido actual de la ventana.

#### `Refresh()`

Refresca la ventana, mostrando los cambios realizados en ésta.

#### `Width() int`

Devuelve el ancho actual de la ventana.

| Valores de retorno | Descripción       |
| ------------------ | ----------------- |
| `int`              | Ancho en píxeles. |

#### `Height() int`

Devuelve la altura actual de la ventana.

| Valores de retorno | Descripción        |
| ------------------ | ------------------ |
| `int`              | Altura en píxeles. |

#### `Resize(newWidth, newHeight int)`

Redimensiona la ventana a un nuevo ancho y alto.

| Argumentos  | Descripción                            |
| ----------- | -------------------------------------- |
| `newWidth`  | Nuevo ancho de la ventana en píxeles.  |
| `newHeight` | Nueva altura de la ventana en píxeles. |

#### `Close()`

Cierra la ventana y libera los recursos utilizados.
