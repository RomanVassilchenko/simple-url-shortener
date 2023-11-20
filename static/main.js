document.addEventListener('DOMContentLoaded', function () {
    const cutButton = document.getElementById('cut');
    const urlInput = document.getElementById('urlInput');
    const aliasInput = document.getElementById('aliasInput');
    const urlOutput = document.getElementById('urlOutput');
    const copyButton = document.getElementById('openModalButton');

    function isValidURL(url) {
        // Regular expression to validate a URL
        const urlRegex = /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([\/\w.-]*)*\/?(\?[^\s]*)?$/;
        return urlRegex.test(url);
    }

    cutButton.addEventListener('click', function () {
        let url = urlInput.value.trim();
        let alias = aliasInput.value.trim()

        if (!url.match(/^https?:\/\//i)) {
            url = "https://" + url;
        }

        if (isValidURL(url)) {
            // Create a JSON object with the URL data
            const requestData = {
                url:url, alias:alias
            };

            // Make the AJAX POST request
            fetch('/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData),
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.status === 'OK' && data.alias) {
                        const currentOrigin = window.location.origin;
                        // Update the output input text with the shortened URL
                        urlOutput.value = `${currentOrigin}/${data.alias}`;
                    } else {
                        alert('Failed to shorten the URL. Error: ' + data.error);
                    }
                })
                .catch((error) => {
                    console.error('Error:', error);
                    alert('An error occurred while shortening the URL. Please try again.');
                });
        } else {
            alert('Please enter a valid URL.');
        }
    });

    // Add event listener to copy button
    copyButton.addEventListener('click', function () {
        // Check if the Clipboard API is available
        if (navigator.clipboard) {
            // Use the Clipboard API to copy the text
            navigator.clipboard.writeText(urlOutput.value)
                .then(function () {
                    // Notify the user that the text has been copied
                    // alert('Copied to clipboard: ' + urlOutput.value);
                })
                .catch(function (err) {
                    console.error('Failed to copy text: ', err);
                });
        } else {
            // Clipboard API is not available, fall back to the old method
            const textArea = document.createElement('textarea');
            textArea.value = urlOutput.value;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            // alert('Copied to clipboard: ' + urlOutput.value);
        }
    });
});