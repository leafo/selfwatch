import React, { memo, useState, useEffect } from 'react';

const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const HOUR_LABELS = [0, 6, 12, 18];

function getColor(p) {
    if (p <= 0) {
        return '#161b22';
    }
    p = Math.min(1, Math.max(0, p));
    return `color-mix(in oklch, #1f2630, #39d353 ${p * 100}%)`;
}

export default memo(function WeeklyHeatmap() {
    const [data, setData] = useState(null);

    useEffect(() => {
        fetch('/api/weekly-heatmap')
            .then(res => res.json())
            .then(setData)
            .catch(console.error);
    }, []);

    if (!data) {
        return (
            <section className="chart-section">
                <h2>Last 7 Days</h2>
                <div className="chart-container loading" />
            </section>
        );
    }

    // Build lookup: { "2024-12-10": { 0: 123, 1: 456, ... } }
    const countByDayHour = {};
    let maxCount = 1;
    data.data.forEach(d => {
        if (!countByDayHour[d.day]) {
            countByDayHour[d.day] = {};
        }
        countByDayHour[d.day][d.hour] = d.count;
        if (d.count > maxCount) maxCount = d.count;
    });

    // Build days array from server-provided date range
    const days = [];
    const startDate = new Date(data.startDate + 'T00:00:00');
    for (let i = 0; i < 7; i++) {
        const date = new Date(startDate);
        date.setDate(startDate.getDate() + i);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        days.push({
            dateStr: `${year}-${month}-${day}`,
            dayName: WEEKDAYS[date.getDay()]
        });
    }

    return (
        <section className="chart-section">
            <h2>Last 7 Days</h2>
            <div className="weekly-heatmap">
                <div className="heatmap-grid">
                    {/* Header row with day labels */}
                    <div className="heatmap-row header-row">
                        <div className="hour-label" />
                        {days.map(({ dateStr, dayName }) => (
                            <div key={dateStr} className="day-label">{dayName}</div>
                        ))}
                    </div>
                    {/* Hour rows */}
                    {Array.from({ length: 24 }, (_, hour) => (
                        <div key={hour} className="heatmap-row">
                            <div className="hour-label">
                                {HOUR_LABELS.includes(hour) ? `${hour}:00` : ''}
                            </div>
                            {days.map(({ dateStr, dayName }) => {
                                const count = countByDayHour[dateStr]?.[hour] || 0;
                                const scaled = count / maxCount;
                                return (
                                    <div
                                        key={dateStr}
                                        className="heatmap-cell"
                                        style={{ backgroundColor: getColor(scaled) }}
                                    >
                                        {count > 0 && (
                                            <div className="heatmap-tooltip">
                                                {dayName} {hour}:00 - {count.toLocaleString()} keys
                                            </div>
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    ))}
                </div>
            </div>
        </section>
    );
});
