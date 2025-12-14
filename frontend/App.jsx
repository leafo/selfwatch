import React, { useState, useEffect } from 'react';
import BarChart from './components/BarChart';
import HourlyActivity from './components/HourlyActivity';
import YearlyActivity from './components/YearlyActivity';

async function fetchData(endpoint) {
    const response = await fetch(endpoint);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${endpoint}`);
    }
    return response.json();
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
            count: countByDay[dayKey] || 0,
            date: dayKey
        });
    }

    return data;
}

export default function App() {
    const [weeklyData, setWeeklyData] = useState(null);
    const [monthlyData, setMonthlyData] = useState(null);
    const [focusedDate, setFocusedDate] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
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
                <HourlyActivity
                    focusedDate={focusedDate}
                    onClearFocus={() => setFocusedDate(null)}
                />

                <section className="chart-section">
                    <h2>Last 7 Days</h2>
                    <div className={`chart-container ${!weeklyData ? 'loading' : ''}`}>
                        {weeklyData && <BarChart data={weeklyData} barClass="weekly-bar" />}
                    </div>
                </section>

                <section className="chart-section">
                    <h2>Last 30 Days</h2>
                    <div className={`chart-container ${!monthlyData ? 'loading' : ''}`}>
                        {monthlyData && (
                            <BarChart
                                data={monthlyData}
                                barClass="monthly-bar"
                                onBarClick={(index, item) => setFocusedDate(item.date)}
                            />
                        )}
                    </div>
                </section>

                <YearlyActivity />
            </main>

            {error && <div style={{ color: 'red', textAlign: 'center' }}>Error: {error}</div>}
        </>
    );
}
