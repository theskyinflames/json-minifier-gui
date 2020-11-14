function dispatch(msg){
    const cmdAbout= 'CMD_ABOUT';
    switch(msg) {
        case cmdAbout:
            alert("(c) theskyinflames, 2021");
            break;
        default:
            document.getElementById('msg').innerHTML = "unknowed command: ".concat(msg);
      }
}

function closeApp(){
    astilectron.sendMessage("close");
}

function minimize(){
    const cmdFormatPrefix = 'CMD_FORMAT_';
    const errPrefix = 'ERR:';
    let raw = document.getElementById('json-raw').value;
    astilectron.sendMessage(cmdFormatPrefix.concat(raw), function(msg) {
        if (msg.startsWith(errPrefix)){
            document.getElementById('msg').innerHTML=msg;
            return;
        }
        document.getElementById('formatted').innerText=msg;
    });
}

function reset(){
    document.getElementById('json-raw').value = '';
    document.getElementById('formatted').innerText = '';
    document.getElementById('msg').innerText = "Please, paste your JSON";
    document.getElementById('json-raw').focus();
}