{{/* vim: set filetype=mustache: */}}

{{/*
Create the name of the service account to use
*/}}
{{- define "statemonitor.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "common.names.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return the proper statemonitor image name
*/}}
{{- define "statemonitor.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "statemonitor.imagePullSecrets" -}}
{{- include "common.images.pullSecrets" (dict "images" (list .Values.image) "global" .Values.global) -}}
{{- end -}}

{{/*
Check if there are rolling tags in the images
*/}}
{{- define "statemonitor.checkRollingTags" -}}
{{- include "common.warnings.rollingTag" .Values.image }}
{{- end -}}
