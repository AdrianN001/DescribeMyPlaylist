const login_button = document.querySelector(".login_button");
window.onload = () => {
    document.body.style.overflow = 'hidden';
}

login_button.addEventListener("click", (e) => {
    e.preventDefault();
    window.location.href = "/login"

})
