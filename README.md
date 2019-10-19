# Artist
This is a simple program that imposes an image over another.
Usage:
        `go run main.go output-name background-name watermark-name watermark-dimensions location`
        e.g. `go run main.go zebra.png sample1.png mark.png 200x200 100x100`
Arguments:
        - output-name: What should the output image be called?
        - background-name: Name of background image.
        - watermark-name: Name of smaller image.
        - watermark-dimensions: Watermark would be resized to this before imposing.
        - location: x & y coordinates to place the watermark.
        
- If smaller image's dimensions are larger than specified, it will be cropped.  