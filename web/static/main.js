document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("myForm");
    const cells = document.querySelectorAll("td[data-id]");
    const header = document.querySelector("header.primary-header");

    // svg map variables
    const svgImage = document.getElementById("svgImage");
    const svgContainer = document.getElementById("svgContainer");
    const vb = svgImage.getAttribute("viewBox").split(" ").map(parseFloat);
    const svgSize = { w: vb[2], h: vb[3] };

    var viewBox = { x: vb[0], y: vb[1], w: vb[2], h: vb[3] };
    var isPanning = false;
    var startPoint = { x: 0, y: 0 };
    var endPoint = { x: 0, y: 0 };
    var scale = 1;

    // set titles to countries on the map
    const countries = document.getElementById("countries").children;
    for (let country of countries) {
        let countryName = country.getAttribute("data-name");
        country.innerHTML = `<title>${countryName}</title>`;
    }

    // handle form submit
    form.onsubmit = e => {

        // reset the previous data, styles and events after submit
        const badVideo = document.getElementById("badVideo");
        if (badVideo) badVideo.remove();

        viewBox = { x: vb[0], y: vb[1], w: vb[2], h: vb[3] };
        isPanning = false;
        startPoint = { x: 0, y: 0 };
        endPoint = { x: 0, y: 0 };
        scale = 1;
        svgContainer.setAttribute("data-scale", scale);
        svgImage.setAttribute('viewBox', `${viewBox.x} ${viewBox.y} ${viewBox.w} ${viewBox.h}`);

        for (let cell of cells) {
            cell.innerText = "";
        }

        for (let country of countries) {
            country.removeAttribute("style");
            country.onmouseenter = () => { };
            country.onmouseleave = () => { };
        }

        e.preventDefault();
        fetch("/", {
            method: "POST",
            body: new FormData(e.target)
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
                    country.onmouseenter = e => { e.target.style.fill = "crimson"; }
                    country.onmouseleave = e => { e.target.style.fill = "red"; }
                }
            });
        });
    };


    svgContainer.onmousewheel = e => {
        e.preventDefault();
        var w = viewBox.w;
        var h = viewBox.h;
        var mx = e.offsetX;
        var my = e.offsetY;
        var dw = w * Math.sign(e.deltaY) * -0.15;
        var dh = h * Math.sign(e.deltaY) * -0.15;
        var dx = dw * mx / svgSize.w;
        var dy = dh * my / svgSize.h;
        viewBox = { x: viewBox.x + dx, y: viewBox.y + dy, w: viewBox.w - dw, h: viewBox.h - dh };
        scale = svgSize.w / viewBox.w;
        svgContainer.setAttribute("data-scale", `${Math.round(scale * 100) / 100}`)
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