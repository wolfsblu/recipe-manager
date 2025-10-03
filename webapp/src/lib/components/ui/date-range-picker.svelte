<script lang="ts">
    import CalendarIcon from "@lucide/svelte/icons/calendar";
    import type { DateRange } from "bits-ui";
    import {
        CalendarDate,
        DateFormatter,
        type DateValue,
        getLocalTimeZone,
        parseDate
    } from "@internationalized/date";
    import { cn } from "$lib/utils.js";
    import { buttonVariants } from "$lib/components/ui/button/index.js";
    import { RangeCalendar } from "$lib/components/ui/range-calendar/index.js";
    import * as Popover from "$lib/components/ui/popover/index.js";

    const df = new DateFormatter("en-US", {
        dateStyle: "medium"
    });

    interface StringDateRange {
        from: string;
        to: string;
    }

    let {
        value = { from: '', to: '' },
        onValueChange,
        placeholder = "Pick a date range"
    }: {
        value?: StringDateRange;
        onValueChange?: (value: StringDateRange) => void;
        placeholder?: string;
    } = $props();

    // Convert string dates to CalendarDate objects
    const stringToCalendarDate = (dateStr: string): CalendarDate | undefined => {
        if (!dateStr) return undefined;
        try {
            return parseDate(dateStr);
        } catch {
            return undefined;
        }
    };

    // Convert CalendarDate to string
    const calendarDateToString = (date: CalendarDate): string => {
        return `${date.year}-${date.month.toString().padStart(2, '0')}-${date.day.toString().padStart(2, '0')}`;
    };

    // Initialize calendar value
    let calendarValue: DateRange = $state({
        start: stringToCalendarDate(value.from),
        end: stringToCalendarDate(value.to)
    });

    // Update calendar value when props change
    $effect(() => {
        calendarValue = {
            start: stringToCalendarDate(value.from),
            end: stringToCalendarDate(value.to)
        };
    });

    // Handle calendar value changes
    const handleCalendarChange = (newValue: DateRange) => {
        calendarValue = newValue;

        if (newValue.start && newValue.end) {
            const stringRange = {
                from: calendarDateToString(newValue.start),
                to: calendarDateToString(newValue.end)
            };
            onValueChange?.(stringRange);
        }
    };

    let startValue: DateValue | undefined = $state(undefined);
</script>

<div class="grid gap-2">
    <Popover.Root>
        <Popover.Trigger
            class={cn(
                buttonVariants({ variant: "outline" }),
                "w-full justify-start text-left font-normal",
                !calendarValue.start && "text-muted-foreground"
            )}
        >
            <CalendarIcon class="mr-2 size-4" />
            {#if calendarValue && calendarValue.start}
                {#if calendarValue.end}
                    {df.format(calendarValue.start.toDate(getLocalTimeZone()))} -
                    {df.format(calendarValue.end.toDate(getLocalTimeZone()))}
                {:else}
                    {df.format(calendarValue.start.toDate(getLocalTimeZone()))}
                {/if}
            {:else if startValue}
                {df.format(startValue.toDate(getLocalTimeZone()))}
            {:else}
                {placeholder}
            {/if}
        </Popover.Trigger>
        <Popover.Content class="w-auto p-0" align="start">
            <RangeCalendar
                bind:value={calendarValue}
                onValueChange={handleCalendarChange}
                class="rounded-md border"
                numberOfMonths={2}
            />
        </Popover.Content>
    </Popover.Root>
</div>