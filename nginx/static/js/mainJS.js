console.log('init');

const token = getMeta("_csrf");
const header = getMeta("_csrf_header");

function getMeta(metaName) {
    const metas = document.getElementsByTagName('meta');

    for (let i = 0; i < metas.length; i++) {
        if (metas[i].getAttribute('name') === metaName) {
            console.log(metas[i].getAttribute('content'));
            return metas[i].getAttribute('content');
        }
    }
    return '';
}

function isValid(s) {
    if (s.length > 50) {
        return false;
    }

    let badSigns = "@#-+$=*^&%<>";
    for (let i = 0; i < badSigns.length; i++) {
        if (s.indexOf(badSigns[i]) > -1) {
            return false;
        }
    }
    return true;
}

document.addEventListener("DOMContentLoaded", function () {
    let input = document.querySelector('.todo-creator_text-input');
    let list = document.querySelector('.todos-list');
    let searchPhrase = document.getElementById('search-phrase');
    initialization();

    function redraw() {
        list.innerHTML = '';
    }


    function initialization() {
        redraw();
    }

    function addItem(filename, count, proximity) {
        list.insertAdjacentHTML(
            "beforeend",
            '<div class="todos-list_item">'
            + '<div class="todos-list_item_text-w">'
            + '<div class="todos-list_item_text" contenteditable="false">'
            + '<span class="todos-toolbar_filters-item">' + filename + '</span>'
            + '<span class="todos-toolbar_filters-item">' + count + '</span>'
            + '<span class="todos-toolbar_filters-item">' + proximity + '</span>'
            + '</div>'
            + '</div>'
            + '</div>'
        );
    }

    function addSearchPhrase(phrase) {
        searchPhrase.innerHTML = phrase
    }

    input.addEventListener("keydown", function (e) {
        if (e.keyCode === 13) {
            e.preventDefault();
            const text = input.value;
            if (text.length > 0) {
                input.value = "";
                if (!isValid(text)) {
                    alert("Input data is not valid");
                    return;
                }

                const formData = new FormData();
                formData.append("search", text);
                const createRequest = new XMLHttpRequest();
                createRequest.open("POST", "http://localhost:80/api");
                createRequest.send(formData);
                createRequest.onreadystatechange = function () {
                    if (createRequest.readyState === XMLHttpRequest.DONE) {
                        // Everything is good, the response was received.
                        if (createRequest.status === 200) { // Perfect!
                            const responseCreate = JSON.parse(createRequest.responseText);
                            console.log(responseCreate);
                            redraw();
                            addSearchPhrase(text);
                            responseCreate.forEach(function (t) {
                                console.log(t.Filename, t.Count, t.Spacing);
                                addItem(t.Filename, t.Count, t.Spacing);
                            })

                        } else {
                            if (createRequest.status === 400) {
                                alert("Incorrect data");
                            }

                        }
                    } else {
                    }
                };
            }
        }
    });
});