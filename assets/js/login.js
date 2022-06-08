let loginBtn = document.getElementById("login_btn");
if (loginBtn){
    loginBtn.onclick = () => {
        let loginData = document.getElementsByClassName("login_input");
        let email = loginData[0].value;
        let password = loginData[1].value;
        if (email !== "" && password !== ""){
            Send("POST", "/user/login", {
                Email: email,
                Password: password
            }, (response)=>{
                if (response){
                    window.location.href="/";
                }
            })
        }
    }
}

