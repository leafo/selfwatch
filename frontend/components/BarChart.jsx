import React from 'react';

const CHART_HEIGHT = 136;

export default function BarChart({ data, barClass }) {
    if (!data || data.length === 0) {
        return (
            <div style={{ color: '#8b949e', textAlign: 'center', width: '100%', alignSelf: 'center' }}>
                No data
            </div>
        );
    }

    const maxCount = Math.max(...data.map(d => d.count));

    return (
        <>
            {data.map((item, index) => {
                const height = maxCount > 0 ? (item.count / maxCount) * CHART_HEIGHT : 0;
                const barHeight = Math.max(height, item.count > 0 ? 2 : 0);

                return (
                    <div key={index} className="bar-wrapper">
                        <div className="bar-area">
                            <div
                                className={`bar ${barClass}`}
                                style={{ height: `${barHeight}px` }}
                            >
                                <div className="bar-tooltip">
                                    {item.count.toLocaleString()} keys
                                </div>
                            </div>
                        </div>
                        <div className="bar-label">{item.label}</div>
                    </div>
                );
            })}
        </>
    );
}
