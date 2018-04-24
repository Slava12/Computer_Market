function Selected(a) {
    var label = a.value;
    var categories = document.getElementsByClassName('category');
    for (i = 0; i < categories.length; i++) {
        var name = categories[i].value
        if (label==name) {
            document.getElementById(name).style.display='block';
        }
        else {
            document.getElementById(name).style.display='none';
        }
    }
}

function SetFeatures() {
    var selectedIndex = document.getElementById("categorySelect").options.selectedIndex;
    var val= document.getElementById("categorySelect").options[selectedIndex].value;

    var features = document.getElementById("features");
    features.value = "";

    var n = document.getElementById(val).children.length;
    for (i = 0; i < n; i++) {
        var elem = document.getElementById(val).children[i];
        var noteI  = elem.lastElementChild;
        if (i == n - 1) {
            features.value += noteI.value;
        }
        else {
            features.value += noteI.value + ";";
        }
    }
}

var imageCount = 0;
function ShowButton() {
    if (imageCount < 10) {
        var div = document.createElement('div');
        div.className = "picture" + imageCount;
        div.innerHTML = "<input class=\"inputFile\" type=\"file\" name=\"file1\" accept=\"image/*\">";

        document.getElementById("pictures").appendChild(div);
        imageCount++;
    }
    else {
        document.getElementById("addPic").style.display = "none";
    }
}