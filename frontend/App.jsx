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

function calculateStats(monthlyData) {
    if (!monthlyData || monthlyData.length < 30) return null;

    const today = monthlyData[29].count;
    const yesterday = monthlyData[28].count;

    const thisWeek = monthlyData.slice(23).reduce((sum, d) => sum + d.count, 0);
    const lastWeek = monthlyData.slice(16, 23).reduce((sum, d) => sum + d.count, 0);

    return {
        today,
        todayDelta: today - yesterday,
        thisWeek,
        weekDelta: thisWeek - lastWeek
    };
}

function formatDelta(delta) {
    if (delta > 0) return `+${delta.toLocaleString()}`;
    return delta.toLocaleString();
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

    const stats = calculateStats(monthlyData);

    return (
        <>
            <header>
                <h1>selfwatch</h1>
                {stats && (
                    <div className="header-stats">
                        <div className="stat-item">
                            <div className="stat-label">Today</div>
                            <div className="stat-value">{stats.today.toLocaleString()}</div>
                            <div className={`stat-delta ${stats.todayDelta >= 0 ? 'positive' : 'negative'}`}>
                                {formatDelta(stats.todayDelta)}
                            </div>
                        </div>
                        <div className="stat-item">
                            <div className="stat-label">This Week</div>
                            <div className="stat-value">{stats.thisWeek.toLocaleString()}</div>
                            <div className={`stat-delta ${stats.weekDelta >= 0 ? 'positive' : 'negative'}`}>
                                {formatDelta(stats.weekDelta)}
                            </div>
                        </div>
                    </div>
                )}
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

            <footer className="build-footer">
                {window.BUILD_INFO?.commit && (
                    <span>
                        {window.BUILD_INFO.commit}
                        {window.BUILD_INFO.date && ` · ${window.BUILD_INFO.date}`}
                    </span>
                )}
            </footer>
        </>
    );
}
