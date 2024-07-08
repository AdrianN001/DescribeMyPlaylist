window.onload = () => {
    document.body.style.overflow = 'hidden';
    // const next_button = document.querySelector(".next_button");

    // next_button.addEventListener("click", () =>{
    //     const url = window.location.href;
    //     const index_of_code = url.indexOf("=")+1;
    //     const playlist_id = url.substring(index_of_code);
    //     window.location.href = `/musical_elements?playlist=${playlist_id}`
    // })


    const background_music_tag = document.getElementById("background_music");
    background_music_tag.volume = 0.3;
}



const ctx = document.getElementById('myChart');
console.log(ctx)

new Chart(ctx, {
  type: 'bar',
  data: {
    labels: ['pop', 'hungarian pop', 'Other', ],
    datasets: [{
      label: '# of songs from a given genre',
      data: [86, 72, 95],
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
