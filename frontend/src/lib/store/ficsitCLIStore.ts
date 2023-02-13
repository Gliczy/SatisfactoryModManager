import { get, readable, writable } from 'svelte/store';
import { cli, ficsitcli_bindings } from '$wailsjs/go/models';
import { AddProfile, CheckForUpdates, DeleteProfile, GetInstallationsInfo, GetInvalidInstalls, GetProfiles, ImportProfile, RenameProfile, SelectInstall, SetProfile } from '$wailsjs/go/ficsitcli_bindings/FicsitCLI';
import { GetFavoriteMods } from '$wailsjs/go/bindings/Settings';
import { readableBinding, writableBinding } from './wailsStoreBindings';
import { tick } from 'svelte';
import { isLaunchingGame } from './generalStore';
import { queue } from 'async';
import { queueAutoStart } from './settingsStore';

export const invalidInstalls = readableBinding<(Error & {path?: string})[]>([], { initialGet: GetInvalidInstalls });

export const installs = readableBinding<ficsitcli_bindings.InstallationInfo[]>([], { initialGet: GetInstallationsInfo });
export const selectedInstall = writable(null as ficsitcli_bindings.InstallationInfo | null);

export const profiles = writableBinding<string[]>([], { initialGet: GetProfiles });
export const selectedProfile = writable(null as string | null);

Promise.all([installs.waitForInit, profiles.waitForInit]).then(() => {
  const i = get(installs);
  if(i.length > 0) {
    selectedInstall.set(get(installs)[0]);
  }
});

selectedInstall.subscribe((i) => {
  const path = i?.info?.path;
  if(path) {
    SelectInstall(path);
    if(i.installation) {
      selectedProfile.set(i.installation.profile);
    }
    checkForUpdates();
  }
});

selectedProfile.subscribe((p) => {
  if(p) {
    SetProfile(p);
    const install = get(selectedInstall);
    if(install && install.installation) {
      install.installation.profile = p;
    }
    checkForUpdates();
  }
});

export async function addProfile(name: string) {
  await AddProfile(name);
  const newProfiles = get(profiles);
  if(!newProfiles.includes(name)) {
    newProfiles.push(name);
    profiles.set(newProfiles);
  }
}

export async function renameProfile(oldName: string, newName: string) {
  await RenameProfile(oldName, newName);
  const newProfiles = get(profiles);
  if(newProfiles.includes(oldName)) {
    const idx = newProfiles.indexOf(oldName);
    newProfiles[idx] = newName;
    profiles.set(newProfiles);
  }
  get(installs).forEach((i) => { if(i.installation.profile === oldName) { i.installation.profile = newName; } });
  if(get(selectedProfile) === oldName) {
    selectedProfile.set(newName);
  }
}

export async function deleteProfile(name: string) {
  await DeleteProfile(name);
  const newProfiles = get(profiles);
  if(newProfiles.includes(name)) {
    const idx = newProfiles.indexOf(name);
    newProfiles.splice(idx, 1);
    profiles.set(newProfiles);
  }
  get(installs).forEach((i) => { if(i.installation.profile === name) { i.installation.profile = 'Default'; } });
  if(get(selectedProfile) === name) {
    selectedProfile.set('Default');
  }
}

export async function importProfile(name: string, filepath: string) {
  await ImportProfile(name, filepath);
  const newProfiles = get(profiles);
  if(!newProfiles.includes(name)) {
    newProfiles.push(name);
    profiles.set(newProfiles);
    tick().then(() => {
      selectedProfile.set(name);
    });
  }
}

export type ProfileMods = { [name: string]: cli.ProfileMod };

export const manifestMods = readableBinding<ProfileMods>({}, { allowNull: false, updateEvent: 'manifestMods' });

export interface LockedMod {
  version: string;
  hash: string;
  link: string;
  dependencies: { [id: string]: string };
}

export type LockFile = { [name: string]: LockedMod };

export const lockfileMods = readableBinding<LockFile>({}, { allowNull: false, updateEvent: 'lockfileMods' });

export interface Progress {
  item: string;
  progress: number;
  message: string;
}

export const progress = readableBinding<Progress | null>(null, { updateEvent: 'progress' });

export const favoriteMods = readableBinding<string[]>([], { updateEvent: 'favoriteMods', initialGet: GetFavoriteMods });

export const isGameRunning = readableBinding(false, { updateEvent: 'isGameRunning', allowNull: false });

export const canModify = readable(true, (set) => {
  const update = () => {
    set(!get(isGameRunning) && !get(progress) && !get(isLaunchingGame));
  };
  const unsubGameRunning = isGameRunning.subscribe(update);
  const unsubProgress = progress.subscribe(update);
  const unsubLaunchingGame = isLaunchingGame.subscribe(update);
  return () => {
    unsubGameRunning();
    unsubProgress();
    unsubLaunchingGame();
  };

});

export const updates = writable<ficsitcli_bindings.Update[]>([]);
export const updateCheckInProgress = writable(false);

export async function checkForUpdates() {
  updateCheckInProgress.set(true);
  const result = await CheckForUpdates();
  updateCheckInProgress.set(false);
  if(result instanceof Error) {
    throw result;
  }
  updates.set(result ?? []);
}

setInterval(checkForUpdates, 1000 * 60 * 5); // Check for updates every 5 minutes

interface QueuedAction {
  mod: string;
  action: 'install' | 'remove' | 'enable' | 'disable';
  func: () => Promise<T>;
}

const queuedActionsInternal = writable<QueuedAction[]>([]);
export const queuedMods = readable<Omit<QueuedAction, 'func'>[]>([], (set) => {
  const unsub = queuedActionsInternal.subscribe((q) => {
    set(q);
  });
  return unsub;
});

const modActionsQueue = queue((task: () => Promise<T>, cb) => {
  const complete = (e: Error | null) => {
    queuedActionsInternal.set(get(queuedActionsInternal).filter((a) => a.func !== task));
    cb(e);
  };
  task().then(complete).catch(complete);
});

modActionsQueue.empty(() => {
  if(!get(queueAutoStart)) {
    modActionsQueue.pause();
  }
});

queueAutoStart.subscribe((val) => {
  if(val) {
    modActionsQueue.resume();
  } else {
    modActionsQueue.pause();
  }
});

export function startQueue() {
  modActionsQueue.resume();
}

export async function addQueuedModAction(mod: string, action: string, func: () => Promise<T>): Promise<T> {
  const queuedAction = { mod, action, func };
  queuedActionsInternal.set([
    ...get(queuedActionsInternal),
    queuedAction,
  ]);
  if(get(queueAutoStart)) {
    startQueue();
  }
  return modActionsQueue.pushAsync(func);
}

export function removeQueuedModAction(mod: string) {
  const queuedAction = get(queuedActionsInternal).find((a) => a.mod === mod);
  if(!queuedAction) {
    return;
  }
  modActionsQueue.remove((a) => a.data === queuedAction.func);
  queuedActionsInternal.set(get(queuedActionsInternal).filter((a) => a.mod !== mod));
}