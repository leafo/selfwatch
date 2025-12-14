import React, { useState, useEffect } from 'react';
import BarChart from './components/BarChart';
import YearlyActivity from './components/YearlyActivity';

async function fetchData(endpoint) {
    const response = await fetch(endpoint);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${endpoint}`);
    }
    return response.json();
}

function formatHourlyData(rawData) {
    const countByHour = {};
    rawData.forEach(d => {
        countByHour[d.Hour] = d.Count;
    });

    const now = new Date();
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

function formatWeeklyData(rawData) {
    const countByDay = {};
    rawData.forEach(d => {
        countByDay[d.Day] = d.Count;
    });

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

    return data;
}

function formatMonthlyData(rawData) {
    const countByDay = {};
    rawData.forEach(d => {
        countByDay[d.Day] = d.Count;
    });

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

    return data;
}

export default function App() {
    const [hourlyData, setHourlyData] = useState(null);
    const [weeklyData, setWeeklyData] = useState(null);
    const [monthlyData, setMonthlyData] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchData('/api/hourly')
            .then(raw => setHourlyData(formatHourlyData(raw)))
            .catch(err => setError(err.message));

        fetchData('/api/daily')
            .then(raw => {
                setWeeklyData(formatWeeklyData(raw));
                setMonthlyData(formatMonthlyData(raw));
            })
            .catch(err => setError(err.message));
    }, []);

    return (
        <>
            <header>
                <h1>selfwatch</h1>
            </header>

            <main>
                <section className="chart-section">
                    <h2>Last 24 Hours</h2>
                    <div className={`chart-container ${!hourlyData ? 'loading' : ''}`}>
                        {hourlyData && <BarChart data={hourlyData} barClass="hourly-bar" />}
                    </div>
                </section>

                <section className="chart-section">
                    <h2>Last 7 Days</h2>
                    <div className={`chart-container ${!weeklyData ? 'loading' : ''}`}>
                        {weeklyData && <BarChart data={weeklyData} barClass="weekly-bar" />}
                    </div>
                </section>

                <section className="chart-section">
                    <h2>Last 30 Days</h2>
                    <div className={`chart-container ${!monthlyData ? 'loading' : ''}`}>
                        {monthlyData && <BarChart data={monthlyData} barClass="monthly-bar" />}
                    </div>
                </section>

                <YearlyActivity />
            </main>

            {error && <div style={{ color: 'red', textAlign: 'center' }}>Error: {error}</div>}
        </>
    );
}
