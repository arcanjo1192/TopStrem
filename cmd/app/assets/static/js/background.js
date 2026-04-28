(function() {
    const bgDiv = document.getElementById('detail-bg');
    if (bgDiv) {
        const bgUrl = bgDiv.getAttribute('data-background');
        if (bgUrl && bgUrl !== '') {
            bgDiv.style.backgroundImage = "url('" + bgUrl + "')";
            bgDiv.style.backgroundSize = "cover";
            bgDiv.style.backgroundPosition = "center 20%";
        }
    }
})();