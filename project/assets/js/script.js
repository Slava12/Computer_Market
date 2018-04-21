var i = 0;
function AddFeature(){
    /*var div = document.createElement('div');
    div.classList.add("elem");
    var sel = document.createElement ('select');
    sel.setAttribute('name', 'feature');*/
    var name1 = "feature" + i;
    alert(name1);
    document.getElementById('last').innerHTML += '<div class="elem">\n' +
        '<select name=name1 aria-required="true">\n' +
        '\t\t\t\t<option disabled selected>Характеристика</option>\n' +
        '\t\t\t\t{{range .}}\n' +
        '\t\t\t\t<option value={{.ID}}>{{.Name}}</option>\n' +
        '\t\t\t\t{{end}}\n' +
        '\t\t\t</select>\n' +
        '\t\t</div>';
    i++;
}