async function fetchData(endpoint) {
    const response = await fetch(endpoint);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${endpoint}`);
    }
    return response.json();
}

function formatNumber(num) {
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'k';
    }
    return num.toString();
}

// Continuous color interpolation for contribution grid
// p: 0 to 1 (0 = lowest activity, 1 = highest activity)
function getColor(p) {
    if (p <= 0) {
        return '#161b22'; // No activity - background color
    }

    // Clamp to 0-1
    p = Math.min(1, Math.max(0, p));

    // Interpolate between dark green and bright green
    const low = { r: 14, g: 68, b: 41 };   // #0e4429
    const high = { r: 57, g: 211, b: 83 }; // #39d353

    const r = Math.round(low.r + (high.r - low.r) * p);
    const g = Math.round(low.g + (high.g - low.g) * p);
    const b = Math.round(low.b + (high.b - low.b) * p);

    return `rgb(${r}, ${g}, ${b})`;
}

function createBarChart(container, data, barClass) {
    container.innerHTML = '';

    if (data.length === 0) {
        container.innerHTML = '<div style="color: #8b949e; text-align: center; width: 100%; align-self: center;">No data</div>';
        return;
    }

    const maxCount = Math.max(...data.map(d => d.count));
    const chartHeight = 136; // pixels available for bars

    data.forEach(item => {
        const wrapper = document.createElement('div');
        wrapper.className = 'bar-wrapper';

        const barArea = document.createElement('div');
        barArea.className = 'bar-area';

        const bar = document.createElement('div');
        bar.className = `bar ${barClass}`;
        const height = maxCount > 0 ? (item.count / maxCount) * chartHeight : 0;
        bar.style.height = Math.max(height, item.count > 0 ? 2 : 0) + 'px';

        const tooltip = document.createElement('div');
        tooltip.className = 'bar-tooltip';
        tooltip.textContent = item.count.toLocaleString() + ' keys';
        bar.appendChild(tooltip);

        barArea.appendChild(bar);

        const label = document.createElement('div');
        label.className = 'bar-label';
        label.textContent = item.label;

        wrapper.appendChild(barArea);
        wrapper.appendChild(label);
        container.appendChild(wrapper);
    });
}

async function renderHourlyChart() {
    const container = document.getElementById('hourly-chart');
    container.classList.add('loading');
    const rawData = await fetchData('/api/hourly');
    container.classList.remove('loading');

    // Create a map of hour -> count from the API data
    const countByHour = {};
    rawData.forEach(d => {
        countByHour[d.Hour] = d.Count;
    });

    // Generate all 24 hours going back from now
    const now = new Date();
    const data = [];

    for (let i = 23; i >= 0; i--) {
        const hourDate = new Date(now);
        hourDate.setHours(now.getHours() - i, 0, 0, 0);

        // Format to match API format: "YYYY-MM-DD HH"
        const year = hourDate.getFullYear();
        const month = String(hourDate.getMonth() + 1).padStart(2, '0');
        const day = String(hourDate.getDate()).padStart(2, '0');
        const hour = String(hourDate.getHours()).padStart(2, '0');
        const hourKey = `${year}-${month}-${day} ${hour}`;

        data.push({
            label: hour + ':00',
            count: countByHour[hourKey] || 0
        });
    }

    createBarChart(container, data, 'hourly-bar');
}

async function renderWeeklyChart() {
    const container = document.getElementById('weekly-chart');
    container.classList.add('loading');
    const rawData = await fetchData('/api/daily');
    container.classList.remove('loading');

    // Create a map of day -> count from the API data
    const countByDay = {};
    rawData.forEach(d => {
        countByDay[d.Day] = d.Count;
    });

    // Generate last 7 days
    const now = new Date();
    const data = [];

    for (let i = 6; i >= 0; i--) {
        const dayDate = new Date(now);
        dayDate.setDate(now.getDate() - i);

        const year = dayDate.getFullYear();
        const month = String(dayDate.getMonth() + 1).padStart(2, '0');
        const day = String(dayDate.getDate()).padStart(2, '0');
        const dayKey = `${year}-${month}-${day}`;

        data.push({
            label: dayDate.toLocaleDateString('en-US', { weekday: 'short' }),
            count: countByDay[dayKey] || 0
        });
    }

    createBarChart(container, data, 'weekly-bar');
}

async function renderMonthlyChart() {
    const container = document.getElementById('monthly-chart');
    container.classList.add('loading');
    const rawData = await fetchData('/api/daily');
    container.classList.remove('loading');

    // Create a map of day -> count from the API data
    const countByDay = {};
    rawData.forEach(d => {
        countByDay[d.Day] = d.Count;
    });

    // Generate last 30 days
    const now = new Date();
    const data = [];

    for (let i = 29; i >= 0; i--) {
        const dayDate = new Date(now);
        dayDate.setDate(now.getDate() - i);

        const year = dayDate.getFullYear();
        const month = String(dayDate.getMonth() + 1).padStart(2, '0');
        const day = String(dayDate.getDate()).padStart(2, '0');
        const dayKey = `${year}-${month}-${day}`;

        data.push({
            label: dayDate.getDate().toString(),
            count: countByDay[dayKey] || 0
        });
    }

    createBarChart(container, data, 'monthly-bar');
}

async function renderContributionGrid() {
    const grid = document.getElementById('contribution-grid');
    grid.classList.add('loading');
    const data = await fetchData('/api/yearly');
    grid.classList.remove('loading');
    grid.innerHTML = '';

    const year = new Date().getFullYear();
    document.getElementById('year-title').textContent = year + ' Activity';

    // Create lookup map
    const countByDay = {};
    data.forEach(d => { countByDay[d.Day] = d.Count; });

    // Calculate max for scaling
    const maxCount = Math.max(...data.map(d => d.Count), 1);

    // Find the first day of the year and what day of week it is
    const startOfYear = new Date(year, 0, 1);
    const startDayOfWeek = startOfYear.getDay(); // 0 = Sunday

    // Calculate the Sunday before or on Jan 1
    const gridStart = new Date(startOfYear);
    gridStart.setDate(gridStart.getDate() - startDayOfWeek);

    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    const weekdays = ['', 'Mon', '', 'Wed', '', 'Fri', ''];
    let currentMonth = -1;

    // Calculate number of weeks to display (53 max)
    const endOfYear = new Date(year, 11, 31);
    const totalDays = Math.ceil((endOfYear - gridStart) / (1000 * 60 * 60 * 24)) + 1;
    const totalWeeks = Math.ceil(totalDays / 7);

    // Create table structure for perfect alignment
    const table = document.createElement('div');
    table.className = 'grid-table';

    // Month labels row
    const monthRow = document.createElement('div');
    monthRow.className = 'grid-row month-row';

    // Empty cell for weekday label column
    const emptyCorner = document.createElement('div');
    emptyCorner.className = 'grid-cell weekday-cell';
    monthRow.appendChild(emptyCorner);

    for (let week = 0; week < Math.min(totalWeeks, 53); week++) {
        const weekStart = new Date(gridStart);
        weekStart.setDate(gridStart.getDate() + (week * 7));
        const monthOfWeek = weekStart.getMonth();

        const monthCell = document.createElement('div');
        monthCell.className = 'grid-cell month-cell';

        if (monthOfWeek !== currentMonth && weekStart.getFullYear() === year) {
            monthCell.textContent = months[monthOfWeek];
            currentMonth = monthOfWeek;
        }
        monthRow.appendChild(monthCell);
    }
    table.appendChild(monthRow);

    // Day rows (Sun through Sat)
    for (let dayOfWeek = 0; dayOfWeek < 7; dayOfWeek++) {
        const dayRow = document.createElement('div');
        dayRow.className = 'grid-row';

        // Weekday label
        const weekdayLabel = document.createElement('div');
        weekdayLabel.className = 'grid-cell weekday-cell';
        weekdayLabel.textContent = weekdays[dayOfWeek];
        dayRow.appendChild(weekdayLabel);

        // Cells for each week
        for (let week = 0; week < Math.min(totalWeeks, 53); week++) {
            const currentDate = new Date(gridStart);
            currentDate.setDate(gridStart.getDate() + (week * 7) + dayOfWeek);

            const cell = document.createElement('div');
            cell.className = 'grid-cell day-cell';

            if (currentDate.getFullYear() === year) {
                const dateStr = currentDate.toISOString().split('T')[0];
                const count = countByDay[dateStr] || 0;

                // Logarithmic scaling for better differentiation across wide ranges
                const scaled = count > 0 ? Math.log(count) / Math.log(maxCount) : 0;
                cell.style.backgroundColor = getColor(scaled);

                cell.title = `${currentDate.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}: ${count.toLocaleString()} keys`;
            } else {
                cell.style.visibility = 'hidden';
            }

            dayRow.appendChild(cell);
        }

        table.appendChild(dayRow);
    }

    grid.appendChild(table);

    // Render legend
    const legend = document.getElementById('contribution-legend');
    legend.innerHTML = '';

    const lessLabel = document.createElement('span');
    lessLabel.textContent = 'Less';
    legend.appendChild(lessLabel);

    const squares = document.createElement('div');
    squares.className = 'legend-squares';

    // Show gradient from 0 to 1 in steps
    const steps = [0, 0.25, 0.5, 0.75, 1.0];
    steps.forEach(p => {
        const square = document.createElement('div');
        square.className = 'day-cell';
        square.style.backgroundColor = getColor(p);
        squares.appendChild(square);
    });

    legend.appendChild(squares);

    const moreLabel = document.createElement('span');
    moreLabel.textContent = 'More';
    legend.appendChild(moreLabel);
}

async function refreshAll() {
    try {
        await Promise.all([
            renderHourlyChart(),
            renderWeeklyChart(),
            renderMonthlyChart(),
            renderContributionGrid()
        ]);
    } catch (error) {
        console.error('Error loading data:', error);
    }
}

// Initial load
document.addEventListener('DOMContentLoaded', refreshAll);
