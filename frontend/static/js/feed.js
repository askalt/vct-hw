let allEvents = [];

async function loadEvents() {
    try {
        const response = await fetch('/api/events');
        if (!response.ok) {
            throw new Error('Failed to load events');
        }

        allEvents = await response.json();
        displayEvents(allEvents);
    } catch (error) {
        console.error('Error loading events:', error);
        eventsList.innerHTML = '<p>Error loading events. Please try again later.</p>';
    }
}

function displayEvents(events) {
    const eventsList = document.getElementById('eventsList');

    if (events.length === 0) {
        eventsList.innerHTML = '<p>No events found.</p>';
        return;
    }

    // Sort events by date (newest first)
    const sortedEvents = events.sort((a, b) => new Date(b.date) - new Date(a.date));

    eventsList.innerHTML = sortedEvents.map(event => `
        <div class="event-card">
            <h3>${escapeHtml(event.name)}</h3>
            <div class="description">${escapeHtml(event.description)}</div>
            <div class="tags">
                ${event.tags.map(tag => `<span class="tag">${escapeHtml(tag)}</span>`).join('')}
            </div>
            <div class="date">${formatDate(event.date)}</div>
        </div>
    `).join('');
}

function applyFilters() {
    const searchTerm = document.getElementById('searchInput').value.toLowerCase();
    const tagFilter = document.getElementById('tagFilter').value.toLowerCase();

    const filteredEvents = allEvents.filter(event => {
        const matchesSearch = !searchTerm ||
            event.description.toLowerCase().includes(searchTerm) ||
            event.name.toLowerCase().includes(searchTerm);

        const matchesTag = !tagFilter ||
            event.tags.some(tag => tag.toLowerCase().includes(tagFilter));

        return matchesSearch && matchesTag;
    });

    displayEvents(filteredEvents);
}

function clearFilters() {
    document.getElementById('searchInput').value = '';
    document.getElementById('tagFilter').value = '';
    displayEvents(allEvents);
}

function escapeHtml(unsafe) {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

// Load events when page loads
document.addEventListener('DOMContentLoaded', loadEvents);

// Add event listeners for real-time filtering
document.getElementById('searchInput').addEventListener('input', applyFilters);
document.getElementById('tagFilter').addEventListener('input', applyFilters);
