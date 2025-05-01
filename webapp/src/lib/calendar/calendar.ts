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

export const isToday = (date: Date) => {
    return date.toDateString() === new Date().toDateString();
}

export const inCurrentMonth = (date: Date) => {
    return date.getMonth() === new Date().getMonth();
}

export const getWeekdays = (locale: string = Settings.defaultLocale) => {
    const weekdays = Info.weekdaysFormat('short', { locale });
    const firstDay = Info.getStartOfWeek({ locale }); // returns 1 for most of Europe, 0 for US
    // Reorder weekdays starting from the locale's first day of the week
    return [...weekdays.slice(firstDay), ...weekdays.slice(0, firstDay)];
}
