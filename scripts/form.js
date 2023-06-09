let map;
let circle;

var cityCenter = { lat: 48.8566, lng: 2.3522 };
const cityRadius = 2000

function getLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(setCityCenter, errorCityCenter, {
            enableHighAccuracy: true,
            timeout: 5000,
            maximumAge: 0
        })
    } else {
        initMap()
    }
}

function setCityCenter(position) {
    cityCenter = {lat: position.coords.latitude, lng: position.coords.longitude}
    initMap()
}

function errorCityCenter(error) {
    initMap()
}

async function initMap() {
    const { Map, Circle } = await google.maps.importLibrary("maps");

    map = new Map(document.getElementById("map"), {
        center: cityCenter,
        zoom: 11,
    });

    circle = new Circle({
        strokeColor: "#FF0000",
        strokeOpacity: 0.8,
        strokeWeight: 2,
        fillColor: "#FF0000",
        fillOpacity: 0.35,
        map,
        center: cityCenter,
        radius: cityRadius,
        editable: true,
        draggable: true
    });  
}

const weekSize = 7;
const daySize = 12;
const hourStart = 9;
const days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

let timeTable = new Array(weekSize);
for (let i = 0; i < weekSize; i++) {
    timeTable[i] = new Array(daySize)
}

let currentPivot = null;

function isChosenTimeTable(i, j) {
    return timeTable[i][j].classList.contains("chosen-time")
}

function setTimeTable(i, j, value) {
    timeTable[i][j].classList.remove("time", "chosen-time")

    if (value) {
        timeTable[i][j].classList.add("chosen-time")
    } else {
        timeTable[i][j].classList.add("time")
    }
}

function updateTimeTable(i, j) {
    value = isChosenTimeTable(i, j)

    if (value) {
        currentPivot = null
        setTimeTable(i, j, false)
    } else {
        if (currentPivot && currentPivot[0] == i) {
            min = Math.min(currentPivot[1], j)
            max = Math.max(currentPivot[1], j)
            for (let k = min; k <= max; k++) {
                setTimeTable(i, k, true)
            }
        } else {
            currentPivot = [i, j]
            setTimeTable(i, j, true)
        }
    }
}

function initTimeTable() {
    const table = document.getElementById("calendar")
    for (let i = 0; i < weekSize; i++) {
        const row = table.insertRow(i)

        const day = row.insertCell(0)
        day.classList.add("day")
        day.innerHTML = "<b>" + days[i] + "</b>"

        for (let j = 1; j < daySize + 1; j++) {
            const hour = row.insertCell(j)
            hour.classList.add("time")
            hour.textContent = (j + hourStart - 1).toString()
            hour.onclick = function() {
                updateTimeTable(i, j - 1)
            }
            timeTable[i][j - 1] = hour
        }
    } 
}

const languages = ["english", "french", "german", "italian", "russian", "spanish"]

const form = document.querySelector("form");
form.addEventListener("submit", (event) => {
    event.preventDefault()

    checkedLanguages = []
    for (let i = 0; i < languages.length; i++) {
        if (document.getElementById(languages[i] + "-checkbox").checked) {
            checkedLanguages.push(languages[i])
        }
    }

    if (checkedLanguages.length == 0) {
        alert("pick at least one language")
        return
    }

    chosenTime = false
    for (let i = 0; i < weekSize; i++) {
        for (let j = 0; j < daySize; j++) {
            if (isChosenTimeTable(i, j)) {
                chosenTime = true
            }
        }
    }

    if (!chosenTime) {
        alert("pick at least one time slot")
        return
    }

    timeTableCondensed = Array(weekSize);
    for (let i = 0; i < weekSize; i++) {
        timeTableCondensed[i] = Array(daySize)
        for (let j = 0; j < daySize; j++) {
            timeTableCondensed[i][j] = isChosenTimeTable(i, j) ? 1 : 0
        }
    }

    const mapInfo = {'lat': circle.center.lat(), 'lng': circle.center.lng(), 'radius': circle.radius}

    var postData = JSON.stringify({
        'name' : document.getElementById("name").value,
        'email': document.getElementById("email").value,
        'contact-info': document.getElementById("contact-info").value,
        'bio': document.getElementById("bio").value,
        'searching-for': document.getElementById("searching-for").value,
        'map': mapInfo,
        'time': timeTableCondensed,
        'lang': checkedLanguages
    });

    var xhr = new XMLHttpRequest()
    xhr.onreadystatechange = () => { 
        if (xhr.readyState == 4) {
            if (xhr.status == 200) {
                location.href = "form-submit.html"
            } else {
                location.href = "error.html"
            }
        }
    }
    xhr.open('POST', 'https://13.38.96.67/submit', true)
    xhr.setRequestHeader('Content-type', 'application/json')
    xhr.send(postData)
});


initTimeTable();
getLocation();
