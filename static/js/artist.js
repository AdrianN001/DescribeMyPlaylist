window.onload = () => {
    document.body.style.overflow = 'hidden';
    const next_button = document.querySelector(".next_button");

    next_button.addEventListener("click", () =>{
        const url = window.location.href;
        const index_of_code = url.indexOf("=")+1;
        const playlist_id = url.substring(index_of_code);
        window.location.href = `/emotional?playlist=${playlist_id}`
    })


    const background_music_tag = document.getElementById("background_music");
    background_music_tag.volume = 0.3;
}



const ctx = document.getElementById('myChart');

new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['Azahriah', 'Rammstein', 'Other Artists'],
    datasets: [{
      label: '# of songs',
      data: [30, 8, 405],
      borderWidth: 1
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
