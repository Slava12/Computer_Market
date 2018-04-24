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

function Set() {
    var selectedIndex = document.getElementById("categorySelect").options.selectedIndex;
    var val= document.getElementById("categorySelect").options[selectedIndex].value;
    console.log(document.getElementById(val));
}