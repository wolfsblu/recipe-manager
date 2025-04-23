<script lang="ts">
    import {Alert, Card, Group, GroupItem, Timeline, TimelineItem} from "flowbite-svelte";
    import {InfoCircleSolid} from "flowbite-svelte-icons";
    import {getDaysInMonth, inCurrentMonth, isToday} from "$lib/calendar/calendar";

    // TODO: Use https://wicky.nillia.ms/cally/components/calendar-range/ for selecting the date range
    // TODO: Display as calendar next to schedule https://tailwindcss.com/plus/ui-blocks/application-ui/data-display/calendars#component-5bb74a1ac1c8cf96e7b3d1ea2bf4570b

    const now = new Date();
    const days = getDaysInMonth(now);

    let groupTimelines = [
        {
            title: '<span class="font-medium text-gray-900 dark:text-white">Jese Leos</span> likes <span class="font-medium text-gray-900 dark:text-white">Bonnie Green\'s</span> post in <span class="font-medium text-gray-900 dark:text-white"> How to start with Flowbite library</span>',
            src: '/images/profile-picture-1.webp',
            alt: 'alt here',
            href: '/',
            isPrivate: true,
            comment: '"I wanted to share a webinar zeroheight."'
        },
        {
            title: '<span class="font-medium text-gray-900 dark:text-white">Bonnie Green</span> react to <span class="font-medium text-gray-900 dark:text-white">Thomas Lean\'s</span> comment',
            src: '/images/profile-picture-2.webp',
            alt: 'alt here',
            href: '/',
            isPrivate: true,
            comment: '"I wanted to share a webinar zeroheight."'
        }
    ]
</script>

<Alert border class="items-start! mb-3" color="blue">
    <span slot="icon">
        <InfoCircleSolid class="w-5 h-5"/>
        <span class="sr-only">Info</span>
    </span>
    <p class="font-medium">
        Work in progress! This page will show a calendar for planning upcoming
        meals.
    </p>
</Alert>

<div class="grid grid-cols-3 gap-3">
    <Card size="none">
        <div class="grid grid-cols-7 overflow-hidden rounded-md border border-gray-200 dark:border-gray-600">
            {#each days as date, i}
                <div class="flex justify-center p-1 border-t border-l border-gray-200 dark:border-gray-600 font-semibold
                {i % 7 === 0 ? 'border-l-0' : ''}
                {i < 7 ? 'border-t-0' : ''}
                {inCurrentMonth(date) ? 'text-gray-700 dark:text-gray-100' : 'bg-gray-100 dark:bg-gray-700'}"
                >
                    <time class="flex items-center justify-center w-10 h-10
                    {isToday(date) ? ' bg-primary-300 rounded-full text-primary-600 dark:text-primary-400 dark:bg-primary-600' : ''}"
                    >
                        {date.getDate()}
                    </time>
                </div>
            {/each}
        </div>
    </Card>
    <Card size="none" class="col-span-2">
        <Timeline>
            <TimelineItem title="Application UI code in Tailwind CSS" date="February 2022">
                <p class="mb-4 text-base font-normal text-gray-500 dark:text-gray-400">Get access to over 20+ pages including a dashboard layout, charts, kanban board, calendar, and pre-order E-commerce & Marketing pages.</p>
            </TimelineItem>
            <TimelineItem title="Application UI code in Tailwind CSS" date="March 2022">
                <p class="text-base font-normal text-gray-500 dark:text-gray-400">All of the pages and components are first designed in Figma and we keep a parity between the two versions even as we update the project.</p>
            </TimelineItem>
            <TimelineItem title="Application UI code in Tailwind CSS" date="April 2022">
                <p class="text-base font-normal text-gray-500 dark:text-gray-400">Get started with dozens of web components and interactive elements built on top of Tailwind CSS.</p>
            </TimelineItem>
        </Timeline>
    </Card>
</div>