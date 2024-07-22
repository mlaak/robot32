function make_modal(reqno){
    document.getElementById('gpt-image-'+reqno).addEventListener('click', () => {
        console.log("clicked");
      document.getElementById('modal-'+reqno).classList.remove('hidden');
    });

    document.getElementById('download-btn').addEventListener('click', () => {
      const link = document.createElement('a');
      link.href = document.getElementById('modal-image-'+reqno).src;
      link.download = 'generated image';
      link.click();
    });

    document.getElementById('animate-btn').addEventListener('click', () => {
      const image = document.getElementById('modal-image-'+reqno);
      image.classList.add('animate-bounce');
      setTimeout(() => image.classList.remove('animate-bounce'), 1000);
    });

    document.getElementById('modal-'+reqno).addEventListener('click', (e) => {
      if (e.target === document.getElementById('modal-'+reqno)) {
        document.getElementById('modal-'+reqno).classList.add('hidden');
      }
    });
}
make_modal(0);
