<template>
    <button class="btn btn-primary"
        @click = "$emit('btn-click')">
        {{ text }}
    </button> 
</template>

<script>
export default {
    name : 'ComponentButton',
    props : {
        text : String,
    },
}
</script>
