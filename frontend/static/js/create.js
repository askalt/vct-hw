document.getElementById('eventForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const formData = new FormData(this);
    const tags = formData.get('tags').split(',').map(tag => tag.trim()).filter(tag => tag);

    const event = {
        name: formData.get('name'),
        description: formData.get('description'),
        tags: tags,
        date: formData.get('date')
    };

    try {
        const response = await fetch('/api/events', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(event)
        });

        const messageDiv = document.getElementById('message');

        if (response.ok) {
            messageDiv.textContent = 'Event created successfully!';
            messageDiv.className = 'message success';
            document.getElementById('eventForm').reset();
        } else {
            const error = await response.text();
            messageDiv.textContent = `Error creating event: ${error}`;
            messageDiv.className = 'message error';
        }
    } catch (error) {
        console.error('Error creating event:', error);
        const messageDiv = document.getElementById('message');
        messageDiv.textContent = 'Error creating event. Please try again.';
        messageDiv.className = 'message error';
    }
});

// Set default date to current datetime
document.addEventListener('DOMContentLoaded', function () {
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('date').value = now.toISOString().slice(0, 16);
});
