# imagen

> A simple image generator with no external dependencies.

> This was created to be run as a Lambda on Zeit. Test it out here: https://imagen-go.br0p0p.now.sh?width=400

![Divider](https://imagen-go.br0p0p.now.sh/?width=882&height=2&color=4c6d4d)

## Example

### Reload this page - the color of this image will change!

![Example Image](https://imagen-go.br0p0p.now.sh/?width=400&height=200)

## API

Set the following params in the query string to control the properties of the generated image.

- `height {integer}` Height of the image. If not set, the width value will be used.
- `width {integer}` Width of the image. If not set, the height value will be used.
- `color {hex color}` The color of the image background. If not set, a random color will be chosen.

