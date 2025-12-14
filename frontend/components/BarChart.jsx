import React, { memo } from 'react';

const CHART_HEIGHT = 136;

export default memo(function BarChart({ data, barClass, onBarClick }) {
    if (!data || data.length === 0) {
        return (
            <div style={{ color: '#8b949e', textAlign: 'center', width: '100%', alignSelf: 'center' }}>
                No data
            </div>
        );
    }

    const maxCount = Math.max(...data.map(d => d.count));
    const clickable = !!onBarClick;

    return (
        <>
            {data.map((item, index) => {
                const height = maxCount > 0 ? (item.count / maxCount) * CHART_HEIGHT : 0;
                const barHeight = Math.max(height, item.count > 0 ? 2 : 0);

                return (
                    <div
                        key={index}
                        className="bar-wrapper"
                        onClick={clickable ? () => onBarClick(index, item) : undefined}
                        style={clickable ? { cursor: 'pointer' } : undefined}
                    >
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
});
