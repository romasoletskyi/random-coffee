const form = document.querySelector("form");
form.addEventListener("submit", (event) => {
    event.preventDefault()

    var postData = JSON.stringify({
        "meet" : document.getElementById("meet").value,
        "satisfaction": document.getElementById("satisfaction").value,
        "add": document.getElementById("add").value
    });

    var xhr = new XMLHttpRequest()
    xhr.onreadystatechange = () => { 
        if (xhr.readyState == 4) {
            if (xhr.status == 200) {
                location.href = "feedback-submit.html"
            } else {
                location.href = "error.html"
            }
        }
    }
    xhr.open('POST', 'https://13.38.96.67/feedback', true)
    xhr.setRequestHeader('Content-type', 'application/json')
    xhr.send(postData)
});