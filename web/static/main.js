document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("myForm");
    const cells = document.querySelectorAll("td[data-id]");
    const header = document.querySelector("header.primary-header");

    var countryElements = document.getElementById("countries").children;
    for (let countryElem of countryElements) {
        let countryName = countryElem.getAttribute("data-name");
        countryElem.innerHTML = `<title>${countryName}</title>`;
    }

    form.addEventListener("submit", event => {

        // reset the previous data, styles and events after submit
        const badVideo = document.getElementById("badVideo");
        if (badVideo) badVideo.remove();

        for (let cell of cells) {
            cell.innerText = "";
        }

        for (let countryElem of countryElements) {
            countryElem.removeAttribute("style");
            countryElem.onmouseenter = () => { };
            countryElem.onmouseleave = () => { };
        }

        event.preventDefault();
        fetch("/", {
            method: "POST",
            body: new FormData(event.target)
        }).then(response => {
            if (!response.ok) {
                const badVideo = document.createElement("span");
                badVideo.setAttribute("id", "badVideo");
                badVideo.innerText = "Not been able to fetch the metadata for this video.";
                if (header) header.appendChild(badVideo);
                return;
            }

            response.json().then(data => {
                for (const [key, value] of Object.entries(data)) {
                    if (key === "regionRestricted" && !value.length) {
                        document.querySelector(`td[data-id=${key}]`).innerText = false;
                        continue
                    } else if (key === "defaultLanguage" && !value) {
                        document.querySelector(`td[data-id=${key}]`).innerText = "none";
                        continue
                    }
                    document.querySelector(`td[data-id=${key}]`).innerText = value;
                }

                for (let countryCode of data.regionRestricted) {
                    let country = document.querySelector(`path[data-id=${countryCode}]`);
                    if (!country) {
                        continue
                    }
                    country.style.fill = "red";
                    country.onmouseenter = event => { event.target.style.fill = "crimson"; }
                    country.onmouseleave = event => { event.target.style.fill = "red"; }
                }
            });
        });
    });

    const svgImage = document.getElementById("svgImage");
    const svgContainer = document.getElementById("svgContainer");
    // const zoomValue = document.getElementById("zoomValue");
    var vb = svgImage.getAttribute("viewBox").split(" ").map(parseFloat);
    var viewBox = { x: vb[0], y: vb[1], w: vb[2], h: vb[3] };
    const svgSize = { w: svgImage.clientWidth, h: svgImage.clientHeight };
    var isPanning = false;
    var startPoint = { x: 0, y: 0 };
    var endPoint = { x: 0, y: 0 };
    var scale = 1;

    svgContainer.onmousewheel = e => {
        e.preventDefault();
        var w = viewBox.w;
        var h = viewBox.h;
        var mx = e.offsetX;
        var my = e.offsetY;
        var dw = w * Math.sign(e.deltaY) * 0.05;
        var dh = h * Math.sign(e.deltaY) * 0.05;
        var dx = dw * mx / svgSize.w;
        var dy = dh * my / svgSize.h;
        viewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w - dw, h: viewBox.h - dh };
        scale = svgSize.w / viewBox.w;
        // zoomValue.innerText = `${Math.round(scale * 100) / 100}`;
        svgImage.setAttribute('viewBox', `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}`);
    }


    svgContainer.onmousedown = e => {
        isPanning = true;
        startPoint = { x: e.x, y: e.y };
    }

    svgContainer.onmousemove = e => {
        if (isPanning) {
            endPoint = { x: e.x, y: e.y };
            var dx = (startPoint.x - endPoint.x) / scale;
            var dy = (startPoint.y - endPoint.y) / scale;
            var movedViewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w, h: viewBox.h };
            svgImage.setAttribute('viewBox', `${movedViewBox.x} ${movedViewBox.y} ${movedViewBox.w} ${movedViewBox.h}`);
        }
    }

    svgContainer.onmouseup = e => {
        if (isPanning) {
            endPoint = { x: e.x, y: e.y };
            var dx = (startPoint.x - endPoint.x) / scale;
            var dy = (startPoint.y - endPoint.y) / scale;
            viewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w, h: viewBox.h };
            svgImage.setAttribute('viewBox', `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}`);
            isPanning = false;
        }
    }

    svgContainer.onmouseleave = () => {
        isPanning = false;
    }
});