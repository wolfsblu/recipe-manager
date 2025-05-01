import {DateTime, Info, Settings} from 'luxon'

export const getDaysInMonth = (date: Date) => {
    const year = date.getFullYear();
    const month = date.getMonth(); // 0 = January

    const firstDayOfMonth = new Date(year, month, 1);
    // Weekday index (Monday = 0, Sunday = 6)
    const startDayIndex = (firstDayOfMonth.getDay() + 6) % 7;

    const totalDays = 42;
    const dates: Date[] = [];

    const startDate = new Date(year, month, 1 - startDayIndex);

    for (let i = 0; i < totalDays; i++) {
        const date = new Date(startDate);
        date.setDate(startDate.getDate() + i);
        dates.push(date);
    }

    return dates;
}

export const isToday = (date: DateTime) => {
    return date.hasSame(DateTime.local(), 'day')
}

export const inCurrentMonth = (date: Date) => {
    return date.getMonth() === new Date().getMonth();
}

export const getWeekdays = (startDate: DateTime, locale = Settings.defaultLocale) => {
    const firstDay = Info.getStartOfWeek({ locale }); // returns 1 for most of Europe, 0 for US
    const weekday = startDate.weekday % 7
    const startOfWeek = startDate.minus({ days: (7 + weekday - firstDay) % 7 })

    return Array.from({ length: 7 }, (_, i) => {
        const day = startOfWeek.plus({ days: i })
        return {
            name: day.toLocaleString({ weekday: 'short' }),
            date: day,
        };
    });
}