document.addEventListener('DOMContentLoaded', function() {
    const shortUrlInput = document.getElementById('shortUrl');
    const copyBtn = document.getElementById('copyBtn');
    const copyText = copyBtn.querySelector('.copy-text');
    const copiedText = copyBtn.querySelector('.copied-text');
    const qrCodeDiv = document.getElementById('qrCode');

    // Генерация QR кода
    if (qrCodeDiv && shortUrlInput) {
        generateQRCode(shortUrlInput.value);
    }

    copyBtn.addEventListener('click', function() {
        shortUrlInput.select();
        shortUrlInput.setSelectionRange(0, 99999);
        
        try {
            navigator.clipboard.writeText(shortUrlInput.value).then(function() {
                showCopyFeedback();
            });
        } catch (err) {
            document.execCommand('copy');
            showCopyFeedback();
        }
    });

    function showCopyFeedback() {
        copyText.style.display = 'none';
        copiedText.style.display = 'inline';
        
        setTimeout(function() {
            copyText.style.display = 'inline';
            copiedText.style.display = 'none';
        }, 2000);
    }

    function generateQRCode(url) {
        // Очищаем контейнер
        qrCodeDiv.innerHTML = '';
        
        // Создаем QR код
        QRCode.toCanvas(url, { 
            width: 200,
            height: 200,
            margin: 1,
            color: {
                dark: '#333333',
                light: '#ffffff'
            }
        }, function(err, canvas) {
            if (err) {
                console.error('Ошибка генерации QR кода:', err);
                qrCodeDiv.innerHTML = '<p>Не удалось сгенерировать QR код</p>';
                return;
            }
            
            qrCodeDiv.appendChild(canvas);
            
            // Добавляем подпись
            const caption = document.createElement('p');
            caption.textContent = 'Отсканируйте QR код';
            caption.style.marginTop = '10px';
            caption.style.color = '#666';
            caption.style.fontSize = '14px';
            qrCodeDiv.appendChild(caption);
        });
    }
});