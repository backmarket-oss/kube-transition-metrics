{{/*
Expand the name of the chart.
*/}}
{{- define "kube-transition-metrics.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kube-transition-metrics.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kube-transition-metrics.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "kube-transition-metrics.labels" -}}
helm.sh/chart: {{ include "kube-transition-metrics.chart" . }}
{{ include "kube-transition-metrics.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "kube-transition-metrics.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kube-transition-metrics.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}


{{/*
Pod labels
*/}}
{{- define "kube-transition-metrics.podLabels" -}}
{{ include "kube-transition-metrics.selectorLabels" . }}
{{ .Values.podLabels | toYaml -}}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "kube-transition-metrics.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "kube-transition-metrics.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the clusterrole
*/}}
{{- define "kube-transition-metrics.clusterRoleName" -}}
{{- if .Values.role.create }}
{{- default (include "kube-transition-metrics.fullname" .) .Values.role.name }}
{{- else }}
{{- default "default" .Values.role.name }}
{{- end }}
{{- end }}

{{/*
Create the name of the clusterrolebinding
*/}}
{{- define "kube-transition-metrics.clusterRoleBindingName" -}}
{{- if .Values.role.create }}
{{- default (include "kube-transition-metrics.fullname" .) .Values.role.binding.name }}
{{- else }}
{{- default "default" .Values.role.binding.name }}
{{- end }}
{{- end }}
