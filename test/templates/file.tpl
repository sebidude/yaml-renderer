Hello world.

this is {{ .data.name }}
I can haz tags: {{ .data.taglist }}
Attribute "name" of obj: {{ .data.obj.name }}
Envvar: {{ .data.obj.env }}
And here is a number: {{ .data.obj.num }}
This var must not be substituted: {{ .data.obj.notset }}
