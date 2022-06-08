function Send(method, uri, data, callback){
    let xhr = new XMLHttpRequest();
    xhr.open(method, uri);
    xhr.onload = function() {
        if (typeof callback === "function"){
            callback(JSON.parse(this.response))
        }
    }
    if (data){
        xhr.send(JSON.stringify(data))
    } else {
        xhr.send()
    }
}

let exit = document.getElementById("exit");
if (exit){
    exit.onclick = () =>{
        Send("DELETE", "/user/exit", null, (response) => {
            if (response){
                window.location.href = "/";
            }
        })
    }
}