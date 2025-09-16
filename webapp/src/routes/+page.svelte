<script lang="ts">
    import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '$lib/components/ui/card';
    import {Button} from '$lib/components/ui/button';
    import * as Chart from '$lib/components/ui/chart';
    import {AreaChart, BarChart, PieChart, Text, type ChartContextValue} from 'layerchart';
    import {curveNatural} from "d3-shape";
    import {scaleUtc, scaleBand} from "d3-scale";
    import {cubicInOut} from "svelte/easing";

    let recipeStats = $state({
        totalRecipes: 123,
        totalUsers: 456,
        totalIngredients: 543,
    });

    const newRecipes = [
        {date: new Date("2024-01-01"), recipes: 186},
        {date: new Date("2024-02-01"), recipes: 305},
        {date: new Date("2024-03-01"), recipes: 237},
        {date: new Date("2024-04-01"), recipes: 73},
        {date: new Date("2024-05-01"), recipes: 209},
        {date: new Date("2024-06-01"), recipes: 214},
    ];

    const popularTags = [
        {name: 'vegetarian', count: 45, color: "var(--color-vegetarian)"},
        {name: 'quick', count: 38, color: "var(--color-quick)"},
        {name: 'healthy', count: 32, color: "var(--color-healthy)"},
        {name: 'comfort', count: 28, color: "var(--color-comfort)"},
        {name: 'dessert', count: 22, color: "var(--color-dessert)"}
    ];

    const cookingTimes = [
        {range: '0-15 min', count: 25},
        {range: '15-30 min', count: 45},
        {range: '30-60 min', count: 35},
        {range: '60+ min', count: 18}
    ];

    const topIngredients = [
        {name: 'salt', count: 89, color: "var(--chart-1)"},
        {name: 'oliveOil', count: 67, color: "var(--chart-2)"},
        {name: 'garlic', count: 54, color: "var(--chart-3)"},
        {name: 'onion', count: 48, color: "var(--chart-4)"},
        {name: 'blackPepper', count: 42, color: "var(--chart-5)"}
    ];

    const chartConfig = {
        recipes: {label: "Recipes", color: "var(--chart-1)"},
        count: {label: "Recipes", color: "var(--chart-2)"},
        vegetarian: {label: "Vegetarian", color: "var(--chart-1)"},
        quick: {label: "Quick", color: "var(--chart-2)"},
        healthy: {label: "Healthy", color: "var(--chart-3)"},
        comfort: {label: "Comfort Food", color: "var(--chart-4)"},
        dessert: {label: "Dessert", color: "var(--chart-5)"},
        cookingTimes: {label: "Cooking Times", color: "var(--chart-2)"},
        salt: {label: 'Salt', color: "var(--chart-3)"},
        oliveOil: {label: 'Olive Oil', color: "var(--chart-4)"},
        garlic: {label: 'Garlic', color: "var(--chart-3)"},
        onion: {label: 'Onion', color: "var(--chart-3)"},
        blackPepper: {label: 'Black Pepper', color: "var(--chart-3)"}
    } satisfies Chart.ChartConfig

    let context = $state<ChartContextValue>();
</script>

<div class="container mx-auto p-6 space-y-6">
    <!-- Header -->
    <div class="text-center space-y-2">
        <h1 class="text-4xl font-bold">Dashboard</h1>
        <p class="text-muted-foreground text-lg">Discover insights from our community</p>
    </div>


    <!-- Stats Overview -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card class="gap-1">
            <CardHeader>
                <CardTitle class="text-sm font-medium">Total Recipes</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="text-2xl font-bold">{recipeStats.totalRecipes}</div>
                <p class="text-xs text-muted-foreground">Shared by our community</p>
            </CardContent>
        </Card>

        <Card class="gap-1">
            <CardHeader>
                <CardTitle class="text-sm font-medium">Active Users</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="text-2xl font-bold">{recipeStats.totalUsers}</div>
                <p class="text-xs text-muted-foreground">Contributing to the platform</p>
            </CardContent>
        </Card>

        <Card class="gap-1">
            <CardHeader>
                <CardTitle class="text-sm font-medium">Unique Ingredients</CardTitle>
            </CardHeader>
            <CardContent>
                <div class="text-2xl font-bold">{recipeStats.totalIngredients}</div>
                <p class="text-xs text-muted-foreground">Available in recipes</p>
            </CardContent>
        </Card>
    </div>

    <!-- CTA Section -->
    <Card class="text-center">
        <CardHeader>
            <CardTitle>Join Our Community</CardTitle>
            <CardDescription>Start sharing your favorite recipes and discover new ones</CardDescription>
        </CardHeader>
        <CardContent>
            <div class="space-x-4">
                <Button href="/auth/register">
                    Get Started
                </Button>
                <Button variant="outline" href="/auth/login">
                    Sign In
                </Button>
            </div>
        </CardContent>
    </Card>

    <!-- Charts Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Recipe Trend -->
        <Card>
            <CardHeader>
                <CardTitle>New Recipes</CardTitle>
                <CardDescription>Recipes added each month</CardDescription>
            </CardHeader>
            <CardContent>
                <Chart.Container config={chartConfig} class="h-[200px] w-full">
                    <AreaChart
                            data={newRecipes}
                            x="date"
                            xScale={scaleUtc()}
                            series={[
                                {
                                    key: "recipes",
                                    label: "Recipes",
                                    color: chartConfig.recipes.color,
                                },
                            ]}
                            axis="x"
                            props={{
                              area: {
                                curve: curveNatural,
                                "fill-opacity": 0.4,
                                line: { class: "stroke-1" },
                                motion: "tween",
                              },
                              xAxis: {
                                format: (v: Date) => v.toLocaleDateString("en-US", { month: "short" }),
                              },
                            }}
                    >
                        {#snippet tooltip()}
                            <Chart.Tooltip
                                    labelFormatter={(v: Date) => v.toLocaleDateString("en-US", { month: "long" })}
                                    indicator="line"
                            />
                        {/snippet}
                    </AreaChart>
                </Chart.Container>
            </CardContent>
        </Card>

        <!-- Popular Recipe Tags -->
        <Card>
            <CardHeader>
                <CardTitle>Popular Tags</CardTitle>
                <CardDescription>Most used recipe tags</CardDescription>
            </CardHeader>
            <CardContent>
                <Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[200px]">
                    <PieChart
                            data={popularTags}
                            key="name"
                            value="count"
                            c="color"
                            innerRadius={50}
                            padding={8}
                            props={{ pie: { motion: "tween" } }}
                    >
                        {#snippet aboveMarks()}
                            <Text
                                    value={String(recipeStats.totalRecipes)}
                                    textAnchor="middle"
                                    verticalAnchor="middle"
                                    class="fill-foreground text-3xl! font-bold"
                                    dy={3}
                            />
                            <Text
                                    value="Tags"
                                    textAnchor="middle"
                                    verticalAnchor="middle"
                                    class="fill-muted-foreground! text-muted-foreground"
                                    dy={22}
                            />
                        {/snippet}
                        {#snippet tooltip()}
                            <Chart.Tooltip hideLabel/>
                        {/snippet}
                    </PieChart>
                </Chart.Container>
            </CardContent>
        </Card>

        <!-- Cooking Time Distribution -->
        <Card>
            <CardHeader>
                <CardTitle>Cooking Time Distribution</CardTitle>
                <CardDescription>Recipes by preparation time</CardDescription>
            </CardHeader>
            <CardContent>
                <Chart.Container config={chartConfig} class="h-[200px] w-full">
                    <BarChart
                            bind:context
                            data={cookingTimes}
                            xScale={scaleBand().padding(0.25)}
                            x="range"
                            axis="x"
                            series={[{ key: "count", label: chartConfig.cookingTimes.label, color: chartConfig.cookingTimes.color}]}
                            props={{
                                bars: {
                                    stroke: "none",
                                    rounded: "all",
                                    radius: 8,
                                    // use the height of the chart to animate the bars
                                    initialY: context?.height,
                                    initialHeight: 0,
                                    motion: {
                                        x: { type: "tween", duration: 500, easing: cubicInOut },
                                        width: { type: "tween", duration: 500, easing: cubicInOut },
                                        height: { type: "tween", duration: 500, easing: cubicInOut },
                                        y: { type: "tween", duration: 500, easing: cubicInOut },
                                    },
                                },
                                highlight: { area: { fill: "none" } },
                            }}
                    >
                        {#snippet tooltip()}
                            <Chart.Tooltip hideLabel/>
                        {/snippet}
                    </BarChart>
                </Chart.Container>
            </CardContent>
        </Card>

        <!-- Most Used Ingredients -->
        <Card>
            <CardHeader>
                <CardTitle>Top Ingredients</CardTitle>
                <CardDescription>Top ingredients across all recipes</CardDescription>
            </CardHeader>
            <CardContent>
                <Chart.Container config={chartConfig} class="h-[200px] w-full">
                    <BarChart
                            bind:context
                            data={topIngredients}
                            x="name"
                            y="count"
                            c="color"
                            cRange={topIngredients.map((c) => c.color)}
                            axis="x"
                            rule={false}
                            props={{
                                bars: {
                                    stroke: "none",
                                    rounded: "all",
                                    radius: 8,
                                    initialY: context?.height,
                                    initialHeight: 0,
                                    motion: {
                                        height: { type: "tween", duration: 500, easing: cubicInOut },
                                        y: { type: "tween", duration: 500, easing: cubicInOut },
                                    },
                                },
                                xAxis: {
                                    format: (d) => chartConfig[d].label,
                                },
                                highlight: { area: { fill: "none" } },
                            }}
                    >
                        {#snippet tooltip()}
                            <Chart.Tooltip hideLabel/>
                        {/snippet}
                    </BarChart>
                </Chart.Container>
            </CardContent>
        </Card>
    </div>
</div>