document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('shortenForm');
    const originalUrlInput = document.getElementById('originalUrl');
    const submitBtn = document.getElementById('submitBtn');
    const btnText = submitBtn.querySelector('.btn-text');
    const btnLoading = submitBtn.querySelector('.btn-loading');
    const errorDiv = document.getElementById('error');
    const resultDiv = document.getElementById('result');
    const shortUrlInput = document.getElementById('shortUrl');
    const copyBtn = document.getElementById('copyBtn');
    const copyText = copyBtn.querySelector('.copy-text');
    const copiedText = copyBtn.querySelector('.copied-text');

    form.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const url = originalUrlInput.value.trim();
        
        if (!isValidUrl(url)) {
            showError('Пожалуйста, введите корректный URL');
            return;
        }

        // Показываем индикатор загрузки
        setLoading(true);
        hideError();

        try {
            const response = await fetch('/shorten', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Произошла ошибка');
            }

            // Показываем результат
            shortUrlInput.value = data.short_url;
            showResult();
            
        } catch (error) {
            showError(error.message);
        } finally {
            setLoading(false);
        }
    });

    copyBtn.addEventListener('click', function() {
        shortUrlInput.select();
        shortUrlInput.setSelectionRange(0, 99999); // Для мобильных устройств
        
        try {
            navigator.clipboard.writeText(shortUrlInput.value).then(function() {
                // Показываем подтверждение копирования
                copyText.style.display = 'none';
                copiedText.style.display = 'inline';
                
                setTimeout(function() {
                    copyText.style.display = 'inline';
                    copiedText.style.display = 'none';
                }, 2000);
            });
        } catch (err) {
            // Fallback для старых браузеров
            document.execCommand('copy');
            copyText.style.display = 'none';
            copiedText.style.display = 'inline';
            
            setTimeout(function() {
                copyText.style.display = 'inline';
                copiedText.style.display = 'none';
            }, 2000);
        }
    });

    function setLoading(loading) {
        if (loading) {
            btnText.style.display = 'none';
            btnLoading.style.display = 'inline';
            submitBtn.disabled = true;
        } else {
            btnText.style.display = 'inline';
            btnLoading.style.display = 'none';
            submitBtn.disabled = false;
        }
    }

    function showError(message) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
    }

    function hideError() {
        errorDiv.style.display = 'none';
    }

    function showResult() {
        resultDiv.style.display = 'block';
        form.style.display = 'none';
    }

    function isValidUrl(string) {
        try {
            new URL(string);
            return true;
        } catch (_) {
            return false;
        }
    }

    // Автофокус на поле ввода
    originalUrlInput.focus();
});