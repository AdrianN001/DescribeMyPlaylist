const now_playing_box = document.getElementById("now_playing_box");
const now_playing_box_src = document.getElementsByClassName("now_playing_src")[0].innerHTML;
now_playing_box.addEventListener('click', (e) =>{
  e.preventDefault();
  window.open(now_playing_box_src, "_blank");
})