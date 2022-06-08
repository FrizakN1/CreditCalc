let scroll = document.getElementById("scroll");
let meaning = document.getElementById("meaning");
if (scroll){
    scroll.oninput = () =>{
        let val = scroll.value;
        scroll.style.background = "-webkit-linear-gradient(left, #ad7e00 0%, #ad7e00 "+val*100/5000000+"%, #202020 "+val*100/5000000+"%, #202020 100%)";
        meaning.value = null;
        meaning.placeholder = val+" ₽";
        calculateSum()
    }
}

if (meaning){
    meaning.onchange = () => {
        if (meaning.value > 8000000){
            meaning.value = 8000000;
        } else if (meaning.value < 30000){
            meaning.value = 30000;
        }
        scroll.value = meaning.value;
        scroll.style.background = "-webkit-linear-gradient(left, #ad7e00 0%, #ad7e00 "+meaning.value*100/5000000+"%, #202020 "+meaning.value*100/5000000+"%, #202020 100%)";
        calculateSum()
    }
}

let scroll2 = document.getElementById("scroll2");
let meaning2 = document.getElementById("meaning2");
if (scroll2){
    scroll2.oninput = () =>{
        let val = scroll2.value;
        scroll2.style.background = "-webkit-linear-gradient(left, #ad7e00 0%, #ad7e00 "+(val-3)*100/57+"%, #202020 "+(val-3)*100/57+"%, #202020 100%)";
        meaning2.value = null;
        meaning2.placeholder = val+" месяца";
        calculateSum()
    }
}

if (meaning2){
    meaning2.onchange = () => {
        if (meaning2.value > 60){
            meaning2.value = 60;
        } else if (meaning2.value < 3){
            meaning2.value = 3;
        }
        scroll2.value = meaning2.value;
        scroll2.style.background = "-webkit-linear-gradient(left, #ad7e00 0%, #ad7e00 "+(meaning2.value-3)*100/57+"%, #202020 "+(meaning2.value-3)*100/57+"%, #202020 100%)";
        calculateSum()
    }
}

function calculateSum(){
    let sum = document.getElementById("sum");
    sum.textContent = `${Math.round(parseInt(scroll.value) * (0.159/12)/(1-1/Math.pow(1+0.159/12,   parseInt(scroll2.value))))} ₽`
}

let applyCredit = document.getElementById("apply_credit");
if (applyCredit){
    applyCredit.onclick = () => {
        Send("PUT", "/api/applyCredit", {
            CreditSum: parseFloat(scroll.value),
            CreditDuration: parseFloat(scroll2.value)
        }, (response) => {
            if (response){
                alert("Кредит успешно оформлен")
            } else {
                window.location.href = "/login"
            }
        })
    }
}