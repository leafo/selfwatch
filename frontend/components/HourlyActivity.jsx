import React, { useState, useEffect, memo } from 'react';
import BarChart from './BarChart';

function formatHourlyData(rawData, offset) {
    const countByHour = {};
    rawData.forEach(d => {
        countByHour[d.Hour] = d.Count;
    });

    const now = new Date();
    now.setDate(now.getDate() - offset);
    const data = [];

    for (let i = 23; i >= 0; i--) {
        const hourDate = new Date(now);
        hourDate.setHours(now.getHours() - i, 0, 0, 0);

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

    return data;
}

function formatHourlyDataForDate(rawData, dateStr) {
    const countByHour = {};
    rawData.forEach(d => {
        countByHour[d.Hour] = d.Count;
    });

    const data = [];
    for (let hour = 0; hour < 24; hour++) {
        const hourStr = String(hour).padStart(2, '0');
        const hourKey = `${dateStr} ${hourStr}`;

        data.push({
            label: hourStr + ':00',
            count: countByHour[hourKey] || 0
        });
    }

    return data;
}

function getDateRange(offset) {
    const now = new Date();
    const end = new Date(now);
    end.setDate(end.getDate() - offset);

    const start = new Date(end);
    start.setHours(start.getHours() - 23);

    const formatDate = (d) => d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

    if (start.toDateString() === end.toDateString()) {
        return formatDate(end);
    }
    return `${formatDate(start)} - ${formatDate(end)}`;
}

function formatFocusedDate(dateStr) {
    const [year, month, day] = dateStr.split('-');
    const date = new Date(year, month - 1, day);
    return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' });
}

export default memo(function HourlyActivity({ focusedDate, onClearFocus }) {
    const [offset, setOffset] = useState(0);
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        setLoading(true);
        setError(null);

        const url = focusedDate
            ? `/api/hourly?date=${focusedDate}`
            : `/api/hourly?offset=${offset}`;

        fetch(url)
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch hourly data');
                return res.json();
            })
            .then(raw => {
                const formatted = focusedDate
                    ? formatHourlyDataForDate(raw, focusedDate)
                    : formatHourlyData(raw, offset);
                setData(formatted);
                setLoading(false);
            })
            .catch(err => {
                setError(err.message);
                setLoading(false);
            });
    }, [offset, focusedDate]);

    const title = focusedDate
        ? formatFocusedDate(focusedDate)
        : `Last 24 Hours · ${getDateRange(offset)}`;

    return (
        <section className="chart-section">
            <div className="section-header">
                <h2>{title}</h2>
                <div className="year-nav">
                    {focusedDate ? (
                        <button className="year-btn" onClick={onClearFocus}>×</button>
                    ) : (
                        <>
                            <button className="year-btn" onClick={() => setOffset(o => o + 1)}>&larr;</button>
                            <button className="year-btn" onClick={() => setOffset(o => o - 1)} disabled={offset <= 0}>&rarr;</button>
                        </>
                    )}
                </div>
            </div>
            <div className={`chart-container ${loading ? 'loading' : ''}`}>
                {data && <BarChart data={data} barClass="hourly-bar" />}
            </div>
            {error && <div style={{ color: 'red' }}>Error: {error}</div>}
        </section>
    );
});
