
const SCREEN_RATIO = window.devicePixelRatio || 1
const SCREEN_WIDTH = window.screen.width * SCREEN_RATIO
const SCREEN_HEIGHT = window.screen.height * SCREEN_RATIO

function _(id) {
    return document.getElementById(id)
}

function PutImagefileOnCanvas(imageFile, canvasId) {
    let image = new Image()
    image.src = URL.createObjectURL(imageFile)

    image.onload = () => {
        let canvas = _(canvasId)
        let ctxt = canvas.getContext("2d")

        let canvasWidth = image.width, canvasHeight = image.height


        if (image.width > image.height) {
            if (image.width > SCREEN_WIDTH) {
                canvasWidth = SCREEN_WIDTH
                canvasHeight = image.height * (canvasWidth / image.width)
            } else {
                canvasWidth = image.width
                canvasHeight = image.height
            }
        } else {
            if (image.height > SCREEN_HEIGHT) {
                canvasHeight = SCREEN_HEIGHT
                canvasWidth = image.width * (canvasHeight / image.height)
            } else {
                canvasWidth = image.width
                canvasHeight = image.height
            }
        }

        canvas.width = canvasWidth
        canvas.height = canvasHeight

        ctxt.drawImage(image, 0, 0, canvasWidth, canvasHeight)
    }
}

function GetImageDatafromCanvas(canvasId) {
    let canvas = _(canvasId)
    let ctxt = canvas.getContext("2d")

    return ctxt.getImageData(0, 0, canvas.width, canvas.height)
}

function Base64ToUint8ClampedArray(base64String) {
    const binaryString = atob(base64String)
    const uint8ClampedArray = new Uint8ClampedArray(binaryString.length)

    for (let i = 0; i < uint8ClampedArray.length; i++) {
        uint8ClampedArray[i] = binaryString.charCodeAt(i)
    }

    return uint8ClampedArray
}

function PutImageDataOnCanvas(canvasId, imageData) {
    let canvas = _(canvasId)
    let ctxt = canvas.getContext("2d")

    canvas.width = imageData.width
    canvas.height = imageData.height
    ctxt.putImageData(imageData, 0, 0)
}

function ShowLoader() {
    _("loader").classList.remove("hide")
}

function HideLoader() {
    _("loader").classList.add("hide")
}

function ShowInitialLoadMessage() {
    let canvas = _("original-image")
    let ctxt = canvas.getContext("2d")

    ctxt.textAlign = "center"
    ctxt.fillText("Click here to select an image", canvas.width/2, canvas.height/2)
}

window.addEventListener("DOMContentLoaded", () => {
    ShowInitialLoadMessage()
    _("original-image").onclick = () => {
        _("image-load-button").click()
    }

    _("image-load-button").onchange = ev => {
        let file = ev.target.files[0]
        PutImagefileOnCanvas(file, "original-image")
    }

    _("grayscale").onclick = async () => {
        ShowLoader()
        const imageData = GetImageDatafromCanvas("original-image")

        const data = { width: imageData.width, height: imageData.height, data: Array.from(imageData.data) }

        const response = await fetch("/oneapi/grayscale", {
            method: "POST",
            body: JSON.stringify(data)
        })
        const responseData = await response.json()
        const greyedImage = new ImageData(Base64ToUint8ClampedArray(responseData.data), Number(responseData.width), Number(responseData.height))

        PutImageDataOnCanvas("oneapi-result-canvas", greyedImage)
        HideLoader()
    }
    
    _("gaussian").onclick = async () => {
        ShowLoader()
        const imageData = GetImageDatafromCanvas("original-image")

        const data = { width: imageData.width, height: imageData.height, data: Array.from(imageData.data) }

        const response = await fetch("/oneapi/gaussian_blur", {
            method: "POST",
            body: JSON.stringify(data)
        })
        const responseData = await response.json()
        const greyedImage = new ImageData(Base64ToUint8ClampedArray(responseData.data), Number(responseData.width), Number(responseData.height))

        PutImageDataOnCanvas("oneapi-result-canvas", greyedImage)
        HideLoader()
    }
})