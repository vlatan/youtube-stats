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
                badVideo.innerText = "Not been able to fetch the info for this video.";
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
});