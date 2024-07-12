
const playlist_id = window.location.href.substring(window.location.href.indexOf("=")+1);


document.getElementById("popularity_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/popularity?playlist=${playlist_id}`
})

document.getElementById("artist_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/artist?playlist=${playlist_id}`
})

document.getElementById("emotion_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/emotional?playlist=${playlist_id}`
})

document.getElementById("musical_element_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/musical_elements?playlist=${playlist_id}`
})

document.getElementById("genre_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/genre?playlist=${playlist_id}`
})


document.getElementById("recents_button").addEventListener('click', (e) =>{
    e.preventDefault()
    window.location.href=`/recents`
})