window.onload = () => {
    const next_button = document.querySelector(".next_button");
    next_button.addEventListener("click", () =>{
        const url = window.location.href;
        const index_of_code = url.indexOf("=")+1;
        const playlist_id = url.substring(index_of_code);
        window.location.href = `/artist?playlist=${playlist_id}`
    })

    const background_music_tag = document.getElementById("background_music");
    background_music_tag.volume = 0.1;


    const now_playing_box = document.getElementById("now_playing_box");
    const now_playing_box_src = document.getElementsByClassName("now_playing_src")[0].innerHTML;
    now_playing_box.addEventListener('click', (e) =>{
      e.preventDefault();
      window.open(now_playing_box_src, "_blank");
    })
}



const ctx = document.getElementById('myChart');

new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['0-20%', '20-40%', '40-60%', '60-80%', '80-100%'],
    datasets: [{
      label: '# of songs in a given "popularity" region',
      data: [12, 19, 3, 5, 2, 3],
      borderWidth: 3
    }]
  },
  options: {
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }
});
