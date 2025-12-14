import React, { memo } from 'react';

const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const WEEKDAYS = ['', 'Mon', '', 'Wed', '', 'Fri', ''];

function getColor(p) {
    if (p <= 0) {
        return '#161b22';
    }
    p = Math.min(1, Math.max(0, p));
    return `color-mix(in oklch, #1f2630, #39d353 ${p * 100}%)`;
}

export default memo(function ContributionGrid({ data, year }) {

    const countByDay = {};
    data.forEach(d => { countByDay[d.Day] = d.Count; });

    const maxCount = Math.max(...data.map(d => d.Count), 1);

    const startOfYear = new Date(year, 0, 1);
    const startDayOfWeek = startOfYear.getDay();

    const gridStart = new Date(startOfYear);
    gridStart.setDate(gridStart.getDate() - startDayOfWeek);

    const endOfYear = new Date(year, 11, 31);
    const totalDays = Math.ceil((endOfYear - gridStart) / (1000 * 60 * 60 * 24)) + 1;
    const totalWeeks = Math.min(Math.ceil(totalDays / 7), 53);

    const monthLabels = [];
    let currentMonth = -1;

    for (let week = 0; week < totalWeeks; week++) {
        const weekStart = new Date(gridStart);
        weekStart.setDate(gridStart.getDate() + (week * 7));
        const monthOfWeek = weekStart.getMonth();

        let label = '';
        if (monthOfWeek !== currentMonth && weekStart.getFullYear() === year) {
            label = MONTHS[monthOfWeek];
            currentMonth = monthOfWeek;
        }
        monthLabels.push(label);
    }

    const rows = [];
    for (let dayOfWeek = 0; dayOfWeek < 7; dayOfWeek++) {
        const cells = [];
        for (let week = 0; week < totalWeeks; week++) {
            const currentDate = new Date(gridStart);
            currentDate.setDate(gridStart.getDate() + (week * 7) + dayOfWeek);

            if (currentDate.getFullYear() === year) {
                const dateStr = currentDate.toISOString().split('T')[0];
                const count = countByDay[dateStr] || 0;
                const scaled = count / maxCount;
                const title = `${currentDate.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}: ${count.toLocaleString()} keys`;

                cells.push(
                    <div
                        key={week}
                        className="grid-cell day-cell"
                        style={{ backgroundColor: getColor(scaled) }}
                        title={title}
                    />
                );
            } else {
                cells.push(
                    <div key={week} className="grid-cell day-cell" style={{ visibility: 'hidden' }} />
                );
            }
        }
        rows.push({ dayOfWeek, cells });
    }

    const legendSteps = [0, 0.25, 0.5, 0.75, 1.0];

    return (
        <>
            <div id="contribution-grid">
                <div className="grid-table">
                    <div className="grid-row month-row">
                        <div className="grid-cell weekday-cell" />
                        {monthLabels.map((label, i) => (
                            <div key={i} className="grid-cell month-cell">{label}</div>
                        ))}
                    </div>
                    {rows.map(({ dayOfWeek, cells }) => (
                        <div key={dayOfWeek} className="grid-row">
                            <div className="grid-cell weekday-cell">{WEEKDAYS[dayOfWeek]}</div>
                            {cells}
                        </div>
                    ))}
                </div>
            </div>
            <div className="legend" id="contribution-legend">
                <span>Less</span>
                <div className="legend-squares">
                    {legendSteps.map((p, i) => (
                        <div key={i} className="day-cell" style={{ backgroundColor: getColor(p) }} />
                    ))}
                </div>
                <span>More</span>
            </div>
        </>
    );
});
