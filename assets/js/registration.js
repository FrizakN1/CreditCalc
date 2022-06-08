let checkEmail = false;
let registrationData = document.getElementsByClassName("registration_input");
if (registrationData.length > 0){
    registrationData[2].onchange = () =>{
        Send("POST", "/user/checkEmail", {
            Email: registrationData[2].value
        }, (response) => {
            if (!response) {
                registrationData[2].style.border = "2px solid red";
                checkEmail = false
            } else {
                registrationData[2].style.border = "2px solid #202020"
                checkEmail = true;
            }
        })
    }
}

let regBtn = document.getElementById("reg_btn");
if (regBtn){
    regBtn.onclick = () =>{
        let name = registrationData[0].value;
        let surname = registrationData[1].value;
        let email = registrationData[2].value;
        let password = registrationData[3].value;
        let rePassword = registrationData[4].value;

        if (name !== "" && surname !== "" && password !== "" && rePassword !== "" && password === rePassword && checkEmail){
            Send("PUT", "/user/addUser", {
                Name: name,
                Surname: surname,
                Email: email,
                Password: password
            }, (response) => {
                if (response){
                    window.location.href = "/login";
                }
            })
        }
    }
}



