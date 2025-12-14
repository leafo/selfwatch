import React, { useState, useEffect } from 'react';
import ContributionGrid from './ContributionGrid';

const currentYear = new Date().getFullYear();

export default function YearlyActivity() {
    const [year, setYear] = useState(currentYear);
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        setLoading(true);
        setError(null);

        fetch(`/api/yearly?year=${year}`)
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch yearly data');
                return res.json();
            })
            .then(raw => {
                setData(raw);
                setLoading(false);
            })
            .catch(err => {
                setError(err.message);
                setLoading(false);
            });
    }, [year]);

    return (
        <section className="contribution-section">
            <div className="section-header">
                <h2>{year} Activity</h2>
                <div className="year-nav">
                    <button className="year-btn" onClick={() => setYear(y => y - 1)}>&larr;</button>
                    <button className="year-btn" onClick={() => setYear(y => y + 1)} disabled={year >= currentYear}>&rarr;</button>
                </div>
            </div>
            <div className={loading ? 'loading' : ''}>
                {data && <ContributionGrid data={data} year={year} />}
            </div>
            {error && <div style={{ color: 'red' }}>Error: {error}</div>}
        </section>
    );
}
